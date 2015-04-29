#include "redis.h"
#include "aof.h"
#include <leveldb/c.h>

void zebraRPush(redisDb *db, struct redisCommand *cmd, robj **argv, int argc) {
  if( cmd->proc == ltrimCommand ||
      cmd->proc == rpushCommand) {
    return;
  }

  if(server.leveldb_state == REDIS_LEVELDB_OFF) {
    return;
  }

  robj *dbq = createStringObject("dbq", 3);
  sds buf = sdsempty();
  buf = catAppendOnlyGenericCommand(buf, argc, argv);
  robj *info = createObject(REDIS_STRING, buf);

  robj *lobj = lookupKeyWrite(db,dbq);
  if (lobj && lobj->type != REDIS_LIST) {
    redisLog(REDIS_WARNING, "db wrong type err: %d", lobj->type);
    return;
  }

  info = tryObjectEncoding(info);
  if (!lobj) {
    lobj = createZiplistObject();
    dbAdd(db,dbq,lobj);
  }
  listTypePush(lobj,info,REDIS_TAIL);

  decrRefCount(dbq);
  decrRefCount(info);
}

int loadLevelDB(char *path) {
  struct redisClient *fakeClient;
  int old_leveldb_state = server.leveldb_state;
  long loops = 0;

  server.leveldb_state = REDIS_LEVELDB_OFF;

  leveldb_t *db;
  leveldb_options_t *options;
  leveldb_readoptions_t *roptions;
  leveldb_iterator_t *iterator;

  char *err = NULL;
  options = leveldb_options_create();
  leveldb_options_set_create_if_missing(options, 1);
  leveldb_options_set_compression(options, 1);
  db = leveldb_open(options, path, &err);

  if (err != NULL) {
    redisLog(REDIS_WARNING, "open leveldb err: %s", err);
    exit(1);
  }

  leveldb_free(err); 
  err = NULL;

  roptions = leveldb_readoptions_create();
  iterator = leveldb_create_iterator(db, roptions);

  fakeClient = createFakeClient();

  char *data = NULL;
  char *value = NULL;
  size_t dataLen = 0;
  size_t valueLen = 0;

  for(leveldb_iter_seek_to_first(iterator); leveldb_iter_valid(iterator); leveldb_iter_next(iterator)) {
    int argc;
    unsigned long len;
    robj **argv;
    struct redisCommand *cmd;

    if (!(loops++ % 1000000)) {
      processEventsWhileBlocked();
      redisLog(REDIS_NOTICE, "load leveldb: %lu", loops);
    }
    data = (char*) leveldb_iter_key(iterator, &dataLen);
    if(data[0] == 'h'){
      argc = 4;
      argv = zmalloc(sizeof(robj*)*argc);
      fakeClient->argc = argc;
      fakeClient->argv = argv;
      argv[0] = createStringObject("hset",4);
      len = data[1];
      argv[1] = createStringObject(data+2,len);
      argv[2] = createStringObject(data+3+len,dataLen-3-len);
      value = (char*) leveldb_iter_value(iterator, &valueLen);
      argv[3] = createStringObject(value, valueLen);
    }else if(data[0] == 's'){
      argc = 3;
      argv = zmalloc(sizeof(robj*)*argc);
      fakeClient->argc = argc;
      fakeClient->argv = argv;
      argv[0] = createStringObject("sadd",4);
      len = data[1];
      argv[1] = createStringObject(data+2,len);
      argv[2] = createStringObject(data+3+len,dataLen-3-len);
    }else if(data[0] == 'z'){
      argc = 4;
      argv = zmalloc(sizeof(robj*)*argc);
      fakeClient->argc = argc;
      fakeClient->argv = argv;
      argv[0] = createStringObject("zadd",4);
      len = data[1];
      argv[1] = createStringObject(data+2,len);
      argv[3] = createStringObject(data+3+len,dataLen-3-len);
      value = (char*) leveldb_iter_value(iterator, &valueLen);
      argv[2] = createStringObject(value, valueLen);
    }else{
      redisLog(REDIS_WARNING,"load leveldb no found type: %d", data[0]);
      continue;
    }
    /* Command lookup */
    cmd = lookupCommand(argv[0]->ptr);
    if (!cmd) {
      redisLog(REDIS_WARNING,"Unknown command '%s' from leveldb", (char*)argv[0]->ptr);
      exit(1);
    }
    /* Run the command in the context of a fake client */
    cmd->proc(fakeClient);

    /* The fake client should not have a reply */
    redisAssert(fakeClient->bufpos == 0 && listLength(fakeClient->reply) == 0);
    /* The fake client should never get blocked */
    redisAssert((fakeClient->flags & REDIS_BLOCKED) == 0);

    /* Clean up. Command code may have changed argv/argc so we use the
     * argv/argc of the client instead of the local variables. */
    freeFakeClientArgv(fakeClient);
  }
  redisLog(REDIS_NOTICE, "load leveldb: %lu", loops);

  freeFakeClient(fakeClient);
  server.leveldb_state = old_leveldb_state;

  leveldb_iter_destroy(iterator);
  leveldb_readoptions_destroy(roptions);
  leveldb_options_destroy(options);
  leveldb_close(db);

  return REDIS_OK;
}


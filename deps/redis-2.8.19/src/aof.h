#ifndef __REDIS_AOF_H
#define __REDIS_AOF_H

#include "redis.h"

sds catAppendOnlyGenericCommand(sds dst, int argc, robj **argv);
struct redisClient *createFakeClient(void);
void freeFakeClientArgv(struct redisClient *c); 
void freeFakeClient(struct redisClient *c); 

#endif


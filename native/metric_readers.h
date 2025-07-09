#ifndef METRIC_READERS_H
#define METRIC_READERS_H

#include <stdint.h>

typedef struct{
    uint64_t mem_total;
    uint64_t mem_free;
    uint64_t mem_available;
    uint64_t buffers;
    uint64_t cached;
} MemoryStatsRaw;

MemoryStatsRaw read_proc_meminfo();

typedef struct{
    uint64_t user;
    uint64_t nice;
    uint64_t system;
    uint64_t idle;
    uint64_t iowait;
    uint64_t irq;
    uint64_t softirq;
} CPUStatsRaw;

CPUStatsRaw read_proc_stat();

#endif
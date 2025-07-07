#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "metric_readers.h"

CPUStatsRaw read_proc_stat() {
    FILE *fp = fopen("/proc/stat", "r");
    CPUStatsRaw stats = {0};

    if(fp == NULL){
        return stats;
    }

    char line[256];
    if(fgets(line, sizeof(line), fp)){
        sscanf(line, "cpu %lu %lu %lu %lu %lu %lu %lu",
                &stats.user,
                &stats.nice,
                &stats.system,
                &stats.idle,
                &stats.iowait,
                &stats.irq,
                &stats.softirq);
    }

    fclose(fp);
    return stats;
}

MemoryStatsRaw read_proc_meminfo(){
    FILE *fp = fopen("/proc/meminfo", "r");
    MemoryStatsRaw stats = {0};

    if(fp == NULL){
        return stats;
    }

    char line[256];
    while(fgets(line, sizeof(line), fp)){
        if (sscanf(line, "MemTotal: %ld kB", &stats.mem_total) == 1) continue;
        if (sscanf(line, "MemFree: %ld kB", &stats.mem_free) == 1) continue;
        if (sscanf(line, "Buffers: %ld kB", &stats.bufffers) == 1) continue;
        if (sscanf(line, "Cached: %ld kB", &stats.cached) == 1) continue;
    }

    fclose(fp);
    return stats;
}
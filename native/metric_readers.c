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
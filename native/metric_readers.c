#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/statfs.h>
#include "metric_readers.h"

int read_network_interfaces(InterfaceStats *stats, int max_interfaces){
    FILE *fp = fopen("/proc/net/dev", "r");
    if(!fp){
        perror("Failed to open /proc/net/dev");
        return -1;
    }

    char line[512];
    int count = 0;

    //skip frst two header lines
    fgets(line, sizeof(line), fp);
    fgets(line, sizeof(line), fp);

    while(fgets(line, sizeof(line), fp) && count < max_interfaces){
        char *colon = strchr(line, ':');
        if (!colon) continue;

        *colon = '\0';
        char iface[32];

        sscanf(line, "%s", iface);
        
    }

    return 0;
}

int read_disk_stats(char *mount, DiskStats *result){
    struct statfs s;
    if(statfs(mount, &s) != 0){
        result->success =0;
        return 1;
    }

    uint64_t total = (s.f_blocks * s.f_bsize) / 1024;
    uint64_t free = (s.f_bfree * s.f_bsize) / 1024;
    uint64_t used = total > free ? (total - free) : 0;

    float usage = 0.0f;
    if(total > 0){
        usage = ((float)used/(float)total) * 100.0f;
    }

    result->total_kb = total;
    result->free_kb = free;
    result->used_kb = used;
    result->usage_pct = usage;
    result->success = 1;

    return 0;
}

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
        if (sscanf(line, "MemAvailable: %lu kB", &stats.mem_available) == 1) continue;
        if (sscanf(line, "Buffers: %ld kB", &stats.buffers) == 1) continue;
        if (sscanf(line, "Cached: %ld kB", &stats.cached) == 1) continue;
    }

    fclose(fp);
    return stats;
}
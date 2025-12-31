//go:build darwin

package cpu

/*
#include <mach/mach.h>
#include <mach/processor_info.h>
#include <mach/mach_host.h>

typedef struct {
    unsigned int user;
    unsigned int system;
    unsigned int idle;
    unsigned int nice;
} cpu_usage_t;

int get_cpu_ticks(cpu_usage_t *usage) {
    mach_msg_type_number_t count;
    processor_cpu_load_info_t cpu_load;
    natural_t num_cpus;

    kern_return_t ret = host_processor_info(
        mach_host_self(),
        PROCESSOR_CPU_LOAD_INFO,
        &num_cpus,
        (processor_info_array_t *)&cpu_load,
        &count
    );

    if (ret != KERN_SUCCESS) {
        return -1;
    }

    usage->user = 0;
    usage->system = 0;
    usage->idle = 0;
    usage->nice = 0;

    for (natural_t i = 0; i < num_cpus; i++) {
        usage->user += cpu_load[i].cpu_ticks[CPU_STATE_USER];
        usage->system += cpu_load[i].cpu_ticks[CPU_STATE_SYSTEM];
        usage->idle += cpu_load[i].cpu_ticks[CPU_STATE_IDLE];
        usage->nice += cpu_load[i].cpu_ticks[CPU_STATE_NICE];
    }

    vm_deallocate(mach_task_self(), (vm_address_t)cpu_load,
        count * sizeof(*cpu_load));

    return 0;
}
*/
import "C"
import (
	"fmt"
	"time"
)

type CPUTicks struct {
	User   uint32
	System uint32
	Idle   uint32
	Nice   uint32
}

func getCPUTicks() (CPUTicks, error) {
	var usage C.cpu_usage_t
	ret := C.get_cpu_ticks(&usage)
	if ret != 0 {
		return CPUTicks{}, fmt.Errorf("failed to get CPU ticks")
	}
	return CPUTicks{
		User:   uint32(usage.user),
		System: uint32(usage.system),
		Idle:   uint32(usage.idle),
		Nice:   uint32(usage.nice),
	}, nil
}

func (t CPUTicks) Total() uint64 {
	return uint64(t.User) + uint64(t.System) + uint64(t.Idle) + uint64(t.Nice)
}

// GetCPUUsage returns CPU usage percentage over the given duration
func GetCpuPercent(interval time.Duration) (float64, error) {
	ticks1, err := getCPUTicks()
	if err != nil {
		return 0, err
	}

	time.Sleep(interval)

	ticks2, err := getCPUTicks()
	if err != nil {
		return 0, err
	}

	totalDiff := float64(ticks2.Total() - ticks1.Total())
	if totalDiff == 0 {
		return 0, nil
	}

	idleDiff := float64(ticks2.Idle - ticks1.Idle)
	usage := 100.0 * (1.0 - idleDiff/totalDiff)

	return usage, nil
}

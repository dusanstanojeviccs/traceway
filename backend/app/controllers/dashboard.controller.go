package controllers

import (
	"backend/app/models"
	"backend/app/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dashboardController struct{}

type DashboardOverviewResponse struct {
	RecentIssues    []models.ExceptionGroup `json:"recentIssues"`
	WorstEndpoints  []models.EndpointStats  `json:"worstEndpoints"`
}

func (d dashboardController) GetDashboardOverview(c *gin.Context) {
	projectId := c.Query("projectId")

	now := time.Now()
	start := now.Add(-24 * time.Hour)

	// Get last 10 issues in the last 24 hours
	recentIssues, _, err := repositories.ExceptionStackTraceRepository.FindGrouped(c, projectId, start, now, 1, 10, "last_seen", "", false)
	if err != nil {
		panic(err)
	}

	// Get 10 worst performing endpoints
	worstEndpoints, err := repositories.TransactionRepository.FindWorstEndpoints(c, projectId, start, now, 10)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, DashboardOverviewResponse{
		RecentIssues:   recentIssues,
		WorstEndpoints: worstEndpoints,
	})
}

func (d dashboardController) GetDashboard(c *gin.Context) {
	projectId := c.Query("projectId")

	now := time.Now()
	var start, end time.Time

	// Parse fromDate parameter
	if fromDateStr := c.Query("fromDate"); fromDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromDateStr); err == nil {
			start = parsed
		} else {
			start = now.Add(-24 * time.Hour)
		}
	} else {
		start = now.Add(-24 * time.Hour)
	}

	// Parse toDate parameter
	if toDateStr := c.Query("toDate"); toDateStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toDateStr); err == nil {
			end = parsed
		} else {
			end = now
		}
	} else {
		end = now
	}

	// Calculate previous period for comparison (same duration before start)
	duration := end.Sub(start)
	prevStart := start.Add(-duration)
	prevEnd := start

	// Calculate aggregation interval based on time range
	intervalMinutes := calculateIntervalMinutes(duration)

	metrics := make([]models.DashboardMetric, 0, 11)

	// 1. Requests count
	requestsTrend, err := repositories.TransactionRepository.CountByInterval(c, projectId, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	requestsCurrent, _ := repositories.TransactionRepository.CountBetween(c, projectId, start, end)
	requestsPrev, _ := repositories.TransactionRepository.CountBetween(c, projectId, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("requests", "Requests", float64(requestsCurrent), "count", requestsTrend, float64(requestsPrev), "requests"))

	// 2. Exceptions count
	exceptionsTrend, err := repositories.ExceptionStackTraceRepository.CountByInterval(c, projectId, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	exceptionsCurrent, _ := repositories.ExceptionStackTraceRepository.CountBetween(c, projectId, start, end)
	exceptionsPrev, _ := repositories.ExceptionStackTraceRepository.CountBetween(c, projectId, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("exceptions", "Exceptions", float64(exceptionsCurrent), "count", exceptionsTrend, float64(exceptionsPrev), "exceptions"))

	// 3. Average Response Time
	avgDurationTrend, err := repositories.TransactionRepository.AvgDurationByInterval(c, projectId, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	avgDurationCurrent := getLastValue(avgDurationTrend)
	avgDurationPrevTrend, _ := repositories.TransactionRepository.AvgDurationByInterval(c, projectId, prevStart, prevEnd, intervalMinutes)
	avgDurationPrev := getAverageValue(avgDurationPrevTrend)
	metrics = append(metrics, buildMetric("avg_response_time", "Avg Response Time", avgDurationCurrent, "ms", avgDurationTrend, avgDurationPrev, "response_time"))

	// 4. Error Rate
	errorRateTrend, err := repositories.TransactionRepository.ErrorRateByInterval(c, projectId, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	errorRateCurrent := getLastValue(errorRateTrend)
	errorRatePrevTrend, _ := repositories.TransactionRepository.ErrorRateByInterval(c, projectId, prevStart, prevEnd, intervalMinutes)
	errorRatePrev := getAverageValue(errorRatePrevTrend)
	metrics = append(metrics, buildMetric("error_rate", "Error Rate", errorRateCurrent, "%", errorRateTrend, errorRatePrev, "error_rate"))

	// 5. CPU Usage
	cpuTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameCpuUsage, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	cpuCurrent := getLastValue(cpuTrend)
	cpuPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameCpuUsage, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("cpu_usage", "CPU Usage", cpuCurrent, "%", cpuTrend, cpuPrev, "cpu"))

	// 6. Memory Usage (MB)
	memTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameMemoryUsage, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	memCurrent := getLastValue(memTrend)
	memPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameMemoryUsage, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("memory_usage", "Memory Usage", memCurrent, "MB", memTrend, memPrev, "memory"))

	// 7. Total System Memory (MB)
	memTotalTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameMemoryTotal, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	memTotalCurrent := getLastValue(memTotalTrend)
	memTotalPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameMemoryTotal, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("memory_total", "Total Memory", memTotalCurrent, "MB", memTotalTrend, memTotalPrev, "memory_total"))

	// 8. Go Routines
	goRoutinesTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameGoRoutines, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	goRoutinesCurrent := getLastValue(goRoutinesTrend)
	goRoutinesPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameGoRoutines, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("go_routines", "Go Routines", goRoutinesCurrent, "", goRoutinesTrend, goRoutinesPrev, "go_routines"))

	// 9. Heap Objects
	heapObjectsTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameHeapObjects, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	heapObjectsCurrent := getLastValue(heapObjectsTrend)
	heapObjectsPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameHeapObjects, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("heap_objects", "Heap Objects", heapObjectsCurrent, "", heapObjectsTrend, heapObjectsPrev, "heap_objects"))

	// 10. Num GC
	numGCTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameNumGC, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	numGCCurrent := getLastValue(numGCTrend)
	numGCPrev, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameNumGC, prevStart, prevEnd)
	metrics = append(metrics, buildMetric("num_gc", "GC Cycles", numGCCurrent, "", numGCTrend, numGCPrev, "num_gc"))

	// 11. GC Pause Total (convert from nanoseconds to milliseconds)
	gcPauseTrend, err := repositories.MetricRecordRepository.GetAverageByInterval(c, projectId, models.MetricNameGCPauseTotal, start, end, intervalMinutes)
	if err != nil {
		panic(err)
	}
	// Convert nanoseconds to milliseconds for display
	for i := range gcPauseTrend {
		gcPauseTrend[i].Value = gcPauseTrend[i].Value / 1_000_000
	}
	gcPauseCurrent := getLastValue(gcPauseTrend)
	gcPausePrevRaw, _ := repositories.MetricRecordRepository.GetAverageBetween(c, projectId, models.MetricNameGCPauseTotal, prevStart, prevEnd)
	gcPausePrev := gcPausePrevRaw / 1_000_000
	metrics = append(metrics, buildMetric("gc_pause", "GC Pause", gcPauseCurrent, "ms", gcPauseTrend, gcPausePrev, "gc_pause"))

	c.JSON(http.StatusOK, models.DashboardResponse{
		Metrics:     metrics,
		LastUpdated: now,
	})
}

func buildMetric(id, name string, current float64, unit string, trend []models.TimeSeriesPoint, prev float64, metricType string) models.DashboardMetric {
	// Convert TimeSeriesPoint to DashboardTrendPoint
	trendPoints := make([]models.DashboardTrendPoint, len(trend))
	for i, p := range trend {
		trendPoints[i] = models.DashboardTrendPoint{
			Timestamp: p.Timestamp,
			Value:     p.Value,
		}
	}

	// Calculate percentage change
	var change24h float64
	if prev > 0 {
		change24h = ((current - prev) / prev) * 100
	}

	// Determine status based on metric type
	status := calculateStatus(current, metricType)

	return models.DashboardMetric{
		ID:        id,
		Name:      name,
		Value:     current,
		Unit:      unit,
		Trend:     trendPoints,
		Change24h: change24h,
		Status:    status,
	}
}

func calculateStatus(value float64, metricType string) string {
	switch metricType {
	case "requests":
		// Lower is worse for requests (less traffic may indicate issues)
		if value < 10 {
			return "critical"
		} else if value < 100 {
			return "warning"
		}
		return "healthy"
	case "exceptions":
		// Higher is worse for exceptions
		if value > 50 {
			return "critical"
		} else if value > 10 {
			return "warning"
		}
		return "healthy"
	case "response_time":
		// Higher is worse for response time
		if value > 500 {
			return "critical"
		} else if value > 200 {
			return "warning"
		}
		return "healthy"
	case "error_rate":
		// Higher is worse for error rate
		if value > 5 {
			return "critical"
		} else if value > 2 {
			return "warning"
		}
		return "healthy"
	case "cpu":
		// Higher is worse for CPU
		if value > 90 {
			return "critical"
		} else if value > 70 {
			return "warning"
		}
		return "healthy"
	case "memory":
		// Higher is worse for memory (MB)
		if value > 900 {
			return "critical"
		} else if value > 700 {
			return "warning"
		}
		return "healthy"
	case "memory_total":
		// Total memory is informational - always healthy
		return "healthy"
	case "go_routines":
		// Higher may indicate goroutine leaks
		if value > 10000 {
			return "critical"
		} else if value > 5000 {
			return "warning"
		}
		return "healthy"
	case "heap_objects":
		// Higher may indicate memory pressure
		if value > 1000000 {
			return "critical"
		} else if value > 500000 {
			return "warning"
		}
		return "healthy"
	case "num_gc", "gc_pause":
		// These are informational, always healthy
		return "healthy"
	}
	return "healthy"
}

func getLastValue(points []models.TimeSeriesPoint) float64 {
	if len(points) == 0 {
		return 0
	}
	return points[len(points)-1].Value
}

func getAverageValue(points []models.TimeSeriesPoint) float64 {
	if len(points) == 0 {
		return 0
	}
	var sum float64
	for _, p := range points {
		sum += p.Value
	}
	return sum / float64(len(points))
}

// calculateIntervalMinutes determines the aggregation bucket size based on the time range duration
func calculateIntervalMinutes(duration time.Duration) int {
	hours := duration.Hours()
	switch {
	case hours < 2:
		return 1 // 1-minute buckets
	case hours < 12:
		return 5 // 5-minute buckets
	case hours < 48:
		return 15 // 15-minute buckets
	case hours < 168: // 7 days
		return 60 // 1-hour buckets
	default:
		return 240 // 4-hour buckets
	}
}

var DashboardController = dashboardController{}

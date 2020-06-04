package execution

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	stub := NewExecutionStub()
	tests := map[string]struct {
		daily, odd, even *Programs
		weekly           *WeeklyPrograms
		temp             *TemperaturePrograms
		err              error
	}{
		"it should return error when is empty": {err: NewInvalidCreateExecution("execution can not be empty")},
		"it should return error when programs are empty": {
			err:    NewInvalidCreateExecution("any program are into execution"),
			daily:  &Programs{},
			weekly: &WeeklyPrograms{},
			odd:    &Programs{},
			even:   &Programs{},
			temp:   &TemperaturePrograms{},
		},
		"it should return execution": {daily: stub.Daily, weekly: stub.Weekly, odd: stub.Odd, even: stub.Even},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := New(tt.daily, tt.weekly, tt.odd, tt.even, tt.temp)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestPrograms_getPrograms(t *testing.T) {
	t.Run("it should return programs in same hour", func(t *testing.T) {
		prgms := Programs{}

		prgm0, _ := NewProgram(15, "22:00", []string{"a", "b"})
		prgm1, _ := NewProgram(20, "15:00", []string{"a", "b"})
		prgm2, _ := NewProgram(30, "15:00", []string{"c", "d"})
		prgm3, _ := NewProgram(20, "21:00", []string{"a", "b"})
		prgms.Add(prgm0)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		prgms.Add(prgm3)

		programsMapped := prgms.GetPrograms()

		programs21 := *programsMapped["21:00"]
		assert.Equal(t, prgm3, programs21[0])
		programs22 := *programsMapped["22:00"]
		assert.Equal(t, prgm0, programs22[0])
		assert.Equal(t, 2, len(*programsMapped["15:00"]))
	})
	t.Run("it should return a single program by hour", func(t *testing.T) {
		prgms := Programs{}

		prgm0, _ := NewProgram(15, "22:00", []string{"a", "b"})
		prgm1, _ := NewProgram(20, "15:00", []string{"a", "b"})
		prgm2, _ := NewProgram(20, "21:00", []string{"a", "b"})
		prgms.Add(prgm0)
		prgms.Add(prgm1)
		prgms.Add(prgm2)

		programsMapped := prgms.GetPrograms()

		programs22 := *programsMapped["22:00"]
		assert.Equal(t, prgm0, programs22[0])
		programs15 := *programsMapped["15:00"]
		assert.Equal(t, prgm1, programs15[0])
		programs21 := *programsMapped["21:00"]
		assert.Equal(t, prgm2, programs21[0])
	})
	t.Run("it should return best program", func(t *testing.T) {
		prgms := Programs{}

		prgm0, _ := NewProgram(15, "15:00", []string{"a", "b"})
		prgm1, _ := NewProgram(20, "15:00", []string{"a", "b"})
		prgm2, _ := NewProgram(18, "15:00", []string{"a", "b"})
		prgms.Add(prgm0)
		prgms.Add(prgm1)
		prgms.Add(prgm2)

		programsMapped := prgms.GetPrograms()
		programs := *programsMapped["15:00"]
		assert.Equal(t, 1, len(programs))
		assert.Equal(t, prgm1.Seconds, programs[0].Seconds)
	})
}

func TestProgram_getBestProgram(t *testing.T) {

	prgr, _ := NewProgram(20, "15:00", []string{"1", "2"})
	new, _ := NewProgram(25, "15:00", []string{"1", "2"})
	new2, _ := NewProgram(15, "15:00", []string{"1", "2"})
	tests := map[string]struct {
		current, new, expected *Program
	}{
		"it should return new object as better": {
			current:  prgr,
			new:      new,
			expected: new,
		},
		"it should return current object as better": {
			current:  prgr,
			new:      new2,
			expected: prgr,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prgm := tt.current.getBestProgram(tt.new)

			assert.Equal(t, tt.expected, prgm)
		})
	}
}

func TestNewProgram(t *testing.T) {
	tests := map[string]struct {
		seconds uint8
		hour    string
		zones   []string
		err     error
	}{
		"it should return error without zones": {
			seconds: 20,
			hour:    "15:00",
			err:     NewInvalidCreateData("zones cannot be empty"),
		},
		"it should return error with empty zones": {
			seconds: 20,
			hour:    "15:00",
			zones:   []string{""},
			err:     NewInvalidCreateData("invalid zone data: "),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			prg, err := NewProgram(tt.seconds, tt.hour, tt.zones)

			assert.Equal(t, tt.err, err)
			if err == nil {
				seconds := float64(tt.seconds)
				assert.Equal(t, seconds, prg.Seconds.Seconds())
			}
		})
	}
}

func TestExecution_GetToday(t *testing.T) {
	t.Run("it should return today daily program", func(t *testing.T) {
		daily := Programs{}
		dailyPrgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		dailyPrgm2, err := NewProgram(40, "16:00", []string{"1", "2"})
		assert.NoError(t, err)
		daily.Add(dailyPrgm1)
		daily.Add(dailyPrgm2)
		week := WeeklyPrograms{}
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*prgms))
		prgm := *prgms
		assert.Equal(t, dailyPrgm1.Seconds, prgm[0].Seconds)
	})

	t.Run("it should return today best daily program in same hour", func(t *testing.T) {
		daily := Programs{}
		dailyPrgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		dailyPrgm2, err := NewProgram(40, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		daily.Add(dailyPrgm1)
		daily.Add(dailyPrgm2)
		week := WeeklyPrograms{}
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*prgms))
		prgm := *prgms
		assert.Equal(t, dailyPrgm2.Seconds, prgm[0].Seconds)
	})

	t.Run("it should return today daily programs in same hour", func(t *testing.T) {
		daily := Programs{}
		dailyPrgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		dailyPrgm2, err := NewProgram(40, "15:00", []string{"3", "4"})
		assert.NoError(t, err)
		daily.Add(dailyPrgm1)
		daily.Add(dailyPrgm2)
		week := WeeklyPrograms{}
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := exec.GetToday(today, 25)

		assert.Equal(t, 2, len(*prgms))
	})

	t.Run("it should return today weekly program", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "16:00", []string{"3", "4"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		weekPrgm := NewWeeklyByDay(&prgms, today.Weekday())
		week.Add(weekPrgm)
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*execPrgms))
	})

	t.Run("it should return today weekly programs", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "15:00", []string{"3", "4"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		weekPrgm := NewWeeklyByDay(&prgms, today.Weekday())
		week.Add(weekPrgm)
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 2, len(*execPrgms))
	})

	t.Run("it should return today best weekly program in same hour", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		weekPrgm := NewWeeklyByDay(&prgms, today.Weekday())
		week.Add(weekPrgm)
		odd := Programs{}
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*execPrgms))
		execPrgm := *execPrgms
		assert.Equal(t, prgm2.Seconds, execPrgm[0].Seconds)
	})

	t.Run("it should return today odd program", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "16:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		odd := Programs{}
		odd.Add(prgm1)
		odd.Add(prgm2)
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*execPrgms))
	})

	t.Run("it should return today odd programs", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "15:00", []string{"3", "4"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		odd := Programs{}
		odd.Add(prgm1)
		odd.Add(prgm2)
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 2, len(*execPrgms))
	})

	t.Run("it should return today best odd program", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		odd := Programs{}
		odd.Add(prgm1)
		odd.Add(prgm2)
		even := Programs{}
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*execPrgms))
		execPrgm := *execPrgms
		assert.Equal(t, prgm2.Seconds, execPrgm[0].Seconds)
	})
	t.Run("it should return today even program", func(t *testing.T) {
		today, err := time.Parse("2006-01-02 15:04:05", "2006-01-03 15:00:00")
		assert.NoError(t, err)

		prgms := Programs{}
		prgm1, err := NewProgram(20, "15:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgm2, err := NewProgram(40, "16:00", []string{"1", "2"})
		assert.NoError(t, err)
		prgms.Add(prgm1)
		prgms.Add(prgm2)
		daily := Programs{}
		week := WeeklyPrograms{}
		odd := Programs{}
		even := Programs{}
		even.Add(prgm1)
		even.Add(prgm2)
		temp := TemperaturePrograms{}
		exec, err := New(&daily, &week, &odd, &even, &temp)
		assert.NoError(t, err)

		execPrgms := exec.GetToday(today, 25)

		assert.Equal(t, 1, len(*execPrgms))
	})
}

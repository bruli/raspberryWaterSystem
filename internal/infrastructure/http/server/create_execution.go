package server

import (
	"net/http"
	"time"

	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
)

type ExecutionBody struct {
	Daily  *Programs
	Weekly *WeeklyPrograms
	Odd    *Programs
	Even   *Programs
	Temp   *TempPrograms
}

func newExecutionBody() *ExecutionBody {
	return &ExecutionBody{}
}

type Program struct {
	Seconds    uint8
	Executions *ExecutionsData
}

func NewProgram(seconds uint8, executions *ExecutionsData) *Program {
	return &Program{Seconds: seconds, Executions: executions}
}

type TempProgram struct {
	Program
	Temperature float32
}

type TempPrograms []*TempProgram

func (tp *TempPrograms) Add(p *TempProgram) {
	*tp = append(*tp, p)
}

type Programs []*Program

func (ex *Programs) Add(e *Program) {
	*ex = append(*ex, e)
}

type ExecutionsData struct {
	Hour  string
	Zones []string
}

func NewExecutionsData(hour string, zones []string) *ExecutionsData {
	return &ExecutionsData{Hour: hour, Zones: zones}
}

type WeeklyPrograms []*Weekly

func (ex *WeeklyPrograms) Add(w *Weekly) {
	*ex = append(*ex, w)
}

type Weekly struct {
	Weekday    time.Weekday
	Executions *Programs
}

func NewWeekly(weekday time.Weekday, executions *Programs) *Weekly {
	return &Weekly{Weekday: weekday, Executions: executions}
}

type createExecution struct {
	create   *execution.Creator
	response *response
	body     *ExecutionBody
}

func (c *createExecution) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := jsoniter.NewDecoder(r.Body)

	b := c.body
	if err := decoder.Decode(b); err != nil {
		c.response.generateJSONErrorResponse(w, err)
		return
	}
	exec, err := c.buildExecution(b)
	if err != nil {
		c.response.generateJSONErrorResponse(w, err)
		return
	}
	if err := c.create.Create(exec); err != nil {
		c.response.generateJSONErrorResponse(w, err)
		return
	}

	c.response.writeJSONResponse(w, http.StatusAccepted, nil)
}

func (c *createExecution) buildExecution(body *ExecutionBody) (execution.Execution, error) {
	daily := execution.Programs{}
	if body.Daily != nil {
		for _, j := range *body.Daily {
			buildProgram, err := c.buildProgram(j)
			if err != nil {
				return execution.Execution{}, err
			}
			daily.Add(buildProgram)
		}
	}
	weekly := execution.WeeklyPrograms{}
	if body.Weekly != nil {
		for _, j := range *body.Weekly {
			buildWeekly, err := c.buildWeekly(j)
			if err != nil {
				return execution.Execution{}, err
			}
			weekly.Add(buildWeekly)
		}
	}
	odd := execution.Programs{}
	if body.Odd != nil {
		for _, j := range *body.Odd {
			buildProgram, err := c.buildProgram(j)
			if err != nil {
				return execution.Execution{}, err
			}
			odd.Add(buildProgram)
		}
	}
	even := execution.Programs{}
	if body.Even != nil {
		for _, j := range *body.Even {
			buildProgram, err := c.buildProgram(j)
			if err != nil {
				return execution.Execution{}, err
			}
			even.Add(buildProgram)
		}
	}
	temp := execution.TemperaturePrograms{}
	if body.Temp != nil {
		for _, j := range *body.Temp {
			buildTemperature, err := c.buildTemperature(j)
			if err != nil {
				return execution.Execution{}, err
			}
			temp.Add(buildTemperature)
		}
	}

	ex, err := execution.New(&daily, &weekly, &odd, &even, &temp)
	if err != nil {
		return execution.Execution{}, err
	}

	return *ex, nil
}

func (c *createExecution) buildProgram(p *Program) (*execution.Program, error) {
	hour := ""
	zones := []string{}
	if p.Executions != nil {
		hour = p.Executions.Hour
		zones = p.Executions.Zones
	}
	return execution.NewProgram(p.Seconds, hour, zones)
}

func (c *createExecution) buildWeekly(w *Weekly) (*execution.Weekly, error) {
	execT := execution.Programs{}
	for _, j := range *w.Executions {
		buildProgram, err := c.buildProgram(j)
		if err != nil {
			return nil, err
		}
		execT.Add(buildProgram)
	}
	return execution.NewWeeklyByDay(&execT, w.Weekday), nil
}

func (c *createExecution) buildTemperature(tp *TempProgram) (*execution.TemperatureProgram, error) {
	return execution.NewTemperatureProgram(tp.Temperature, tp.Seconds, tp.Executions.Hour, tp.Executions.Zones)
}

func newCreateExecution(create *execution.Creator, logger logger.Logger) *createExecution {
	return &createExecution{create: create, response: newResponse(logger), body: newExecutionBody()}
}

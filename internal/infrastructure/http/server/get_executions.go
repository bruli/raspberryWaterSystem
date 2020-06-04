package server

import (
	"net/http"

	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/logger"
	jsoniter "github.com/json-iterator/go"
)

type getExecutions struct {
	getter   *execution.Getter
	response *response
}

func newGetExecutions(read *execution.Getter, log logger.Logger) *getExecutions {
	return &getExecutions{getter: read, response: newResponse(log)}
}

func (g *getExecutions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exec, err := g.getter.Get()
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}

	body, err := jsoniter.Marshal(g.buildExecutionBody(exec))
	if err != nil {
		g.response.generateJSONErrorResponse(w, err)
		return
	}

	g.response.writeJSONResponse(w, http.StatusOK, body)
}

func (g *getExecutions) buildExecutionBody(exec *execution.Execution) ExecutionBody {
	weekly := WeeklyPrograms{}
	if exec.Weekly != nil {
		for _, prgms := range *exec.Weekly {
			weekly.Add(g.buildWeekly(prgms))
		}
	}
	var daily, odd, even *Programs
	if exec.Daily != nil {
		daily = g.buildPrograms(exec.Daily)
	}
	if exec.Odd != nil {
		odd = g.buildPrograms(exec.Odd)
	}
	if exec.Even != nil {
		even = g.buildPrograms(exec.Even)
	}
	temp := TempPrograms{}
	if exec.Temp != nil {
		for _, prgms := range *exec.Temp {
			pgr := g.buildProgram(&prgms.Program)
			tempProgram := TempProgram{Program: *pgr, Temperature: prgms.Temperature}
			temp.Add(&tempProgram)
		}
	}
	return ExecutionBody{
		Daily:  daily,
		Weekly: &weekly,
		Odd:    odd,
		Even:   even,
		Temp:   &temp,
	}
}

func (g *getExecutions) buildPrograms(p *execution.Programs) *Programs {
	programs := Programs{}
	for _, program := range *p {
		programs.Add(g.buildProgram(program))
	}
	return &programs
}

func (g *getExecutions) buildProgram(p *execution.Program) *Program {
	return NewProgram(
		uint8((p.Seconds.Seconds())),
		NewExecutionsData(p.Executions.Hour.Format("15:04"), p.Executions.Zones),
	)
}

func (g *getExecutions) buildWeekly(p *execution.Weekly) *Weekly {
	return NewWeekly(p.Weekday, g.buildPrograms(p.Executions))
}

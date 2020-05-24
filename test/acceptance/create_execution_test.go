package acceptance

import (
	"github.com/bruli/raspberryWaterSystem/internal/execution"
	"github.com/bruli/raspberryWaterSystem/internal/infrastructure/http/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateAndGetExecution(t *testing.T) {
	exec := getExecutionBody()
	body, err := jsoniter.Marshal(&exec)
	assert.NoError(t, err)

	const endpoint = "/executions"
	resp, err := sendRequest(http.MethodPut, endpoint, body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	resp, err = sendRequest(http.MethodGet, endpoint, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody := server.ExecutionBody{}
	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = jsoniter.Unmarshal(data, &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, len(*exec.Daily), len(*responseBody.Daily))

}

func getExecutionBody() server.ExecutionBody {
	stub := execution.NewExecutionStub()
	daily := server.Programs{}
	weeklyPrgms := server.WeeklyPrograms{}
	odd := server.Programs{}
	even := server.Programs{}
	for _, dailyPrg := range *stub.Daily {
		pr := getProgram(dailyPrg)
		daily.Add(&pr)

	}
	for _, weeklyPrg := range *stub.Weekly {
		prgms := server.Programs{}
		for _, exec := range *weeklyPrg.Executions {
			pr := getProgram(exec)
			prgms.Add(&pr)
		}
		weeklyPrgms.Add(&server.Weekly{Weekday: 0, Executions: &prgms})
	}
	for _, odPrg := range *stub.Odd {
		pr := getProgram(odPrg)
		odd.Add(&pr)

	}
	for _, evenPrg := range *stub.Even {
		pr := getProgram(evenPrg)
		even.Add(&pr)

	}
	body := server.ExecutionBody{
		Daily:  &daily,
		Weekly: &weeklyPrgms,
		Odd:    &odd,
		Even:   &even,
	}

	return body
}

func getProgram(prg *execution.Program) server.Program {
	data := server.ExecutionsData{Hour: prg.Executions.Hour.Format("15:04"), Zones: prg.Executions.Zones}
	return server.Program{Seconds: uint8(prg.Seconds.Seconds()), Executions: &data}
}

package llm_caller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"llm_hub/conf"
	"testing"
)

func Test_gptLLMCaller_Call(t *testing.T) {
	conf.Init()
	ctx := context.Background()
	caller, err := NewGptLLMCaller(ctx, "hello world", 0.5, 128)
	assert.Nil(t, err)
	completion, err := caller.Call(ctx, "hello world")
	assert.Nil(t, err)
	t.Logf("gpt call, completion = %v", completion)
}

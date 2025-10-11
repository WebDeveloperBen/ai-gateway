package policies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/WebDeveloperBen/ai-gateway/internal/logger"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/google/cel-go/cel"
)

// CELPolicy represents a policy that uses CEL expressions for evaluation
type CELPolicy struct {
	policyType    model.PolicyType
	preCheckExpr  cel.Program // Compiled CEL expression for pre-check
	postCheckExpr cel.Program // Compiled CEL expression for post-check (optional)
	config        []byte      // Original config for reference
}

// NewCELPolicy creates a new CEL-based policy
func NewCELPolicy(policyType model.PolicyType, config []byte) (*CELPolicy, error) {
	// Parse config to get CEL expressions
	type celConfig struct {
		PreCheckExpression  string `json:"pre_check_expression"`
		PostCheckExpression string `json:"post_check_expression,omitempty"`
	}

	var cfg celConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, fmt.Errorf("invalid CEL policy config: %w", err)
	}

	// Create CEL environment with available variables
	env, err := cel.NewEnv(
		cel.Variable("request_size_bytes", cel.IntType),
		cel.Variable("estimated_tokens", cel.IntType),
		cel.Variable("model", cel.StringType),
		cel.Variable("org_id", cel.StringType),
		cel.Variable("app_id", cel.StringType),
		cel.Variable("prompt_tokens", cel.IntType),
		cel.Variable("completion_tokens", cel.IntType),
		cel.Variable("total_tokens", cel.IntType),
		cel.Variable("latency_ms", cel.IntType),
		cel.Variable("response_size_bytes", cel.IntType),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Compile pre-check expression
	var preCheck cel.Program
	if cfg.PreCheckExpression != "" {
		ast, issues := env.Compile(cfg.PreCheckExpression)
		if issues != nil && issues.Err() != nil {
			return nil, fmt.Errorf("failed to compile pre-check expression: %w", issues.Err())
		}
		preCheck, err = env.Program(ast)
		if err != nil {
			return nil, fmt.Errorf("failed to create pre-check program: %w", err)
		}
	}

	// Compile post-check expression (optional)
	var postCheck cel.Program
	if cfg.PostCheckExpression != "" {
		ast, issues := env.Compile(cfg.PostCheckExpression)
		if issues != nil && issues.Err() != nil {
			return nil, fmt.Errorf("failed to compile post-check expression: %w", issues.Err())
		}
		postCheck, err = env.Program(ast)
		if err != nil {
			return nil, fmt.Errorf("failed to create post-check program: %w", err)
		}
	}

	return &CELPolicy{
		policyType:    policyType,
		preCheckExpr:  preCheck,
		postCheckExpr: postCheck,
		config:        config,
	}, nil
}

// Type returns the policy type
func (p *CELPolicy) Type() model.PolicyType {
	return p.policyType
}

// PreCheck evaluates the CEL pre-check expression
func (p *CELPolicy) PreCheck(ctx context.Context, req *PreRequestContext) error {
	if p.preCheckExpr == nil {
		return nil // No pre-check expression
	}

	// Build evaluation context
	vars := map[string]any{
		"request_size_bytes": req.RequestSizeBytes,
		"estimated_tokens":   req.EstimatedTokens,
		"model":              req.Model,
		"org_id":             req.OrgID,
		"app_id":             req.AppID,
	}

	// Evaluate expression
	out, _, err := p.preCheckExpr.Eval(vars)
	if err != nil {
		return fmt.Errorf("CEL evaluation error: %w", err)
	}

	// Expression should return a boolean
	result, ok := out.Value().(bool)
	if !ok {
		return fmt.Errorf("CEL expression must return boolean, got %T", out.Value())
	}

	if !result {
		return fmt.Errorf("policy %s failed: CEL expression returned false", p.policyType)
	}

	return nil
}

// PostCheck evaluates the CEL post-check expression (async)
func (p *CELPolicy) PostCheck(ctx context.Context, req *PostRequestContext) {
	if p.postCheckExpr == nil {
		return // No post-check expression
	}

	// Build evaluation context
	vars := map[string]any{
		"prompt_tokens":       req.ActualTokens.PromptTokens,
		"completion_tokens":   req.ActualTokens.CompletionTokens,
		"total_tokens":        req.ActualTokens.TotalTokens,
		"latency_ms":          req.LatencyMs,
		"response_size_bytes": req.ResponseSizeBytes,
		"model":               req.ModelName,
		"org_id":              req.OrgID,
		"app_id":              req.AppID,
	}

	// Evaluate expression (errors are logged but don't block)
	_, _, err := p.postCheckExpr.Eval(vars)
	if err != nil {
		logger.GetLogger(ctx).Error().
			Err(err).
			Str("app_id", req.AppID).
			Str("policy_type", string(p.policyType)).
			Msg("CEL post-check evaluation failed")
		return
	}

	// Post-check is for logging/metrics, not enforcement
}

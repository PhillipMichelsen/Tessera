package node

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Task definition
type Task struct {
	Instructions []Instruction `yaml:"instructions"`
}

// Instruction struct that holds type and raw YAML args
type Instruction struct {
	Type string `yaml:"type"`
	Args any    `yaml:"args"`
}

type CreateWorkerInstructionArgs struct {
	WorkerType string    `yaml:"worker_type"`
	WorkerUUID uuid.UUID `yaml:"worker_uuid"`
}

type RemoveWorkerInstructionArgs struct {
	WorkerUUID uuid.UUID `yaml:"worker_uuid"`
}

type StartWorkerInstructionArgs struct {
	WorkerUUID      uuid.UUID `yaml:"worker_uuid"`
	WorkerRawConfig []byte    `yaml:"worker_raw_config"`
}

type StopWorkerInstructionArgs struct {
	WorkerUUID uuid.UUID `yaml:"worker_uuid"`
}

func ParseTaskFromYaml(yamlBytes []byte) (Task, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return Task{}, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// Ensure we have a document node.
	if len(root.Content) == 0 {
		return Task{}, fmt.Errorf("empty yaml document")
	}
	doc := root.Content[0]

	// Find the "instructions" node.
	var instructionsNode *yaml.Node
	for i := 0; i < len(doc.Content); i += 2 {
		keyNode := doc.Content[i]
		if keyNode.Value == "instructions" {
			instructionsNode = doc.Content[i+1]
			break
		}
	}
	if instructionsNode == nil {
		return Task{}, fmt.Errorf("instructions not found")
	}
	if instructionsNode.Kind != yaml.SequenceNode {
		return Task{}, fmt.Errorf("instructions is not a sequence")
	}

	var task Task
	// Process each instruction.
	for _, instrNode := range instructionsNode.Content {
		var instrType string
		var argsNode *yaml.Node
		// Each instruction is a mapping node.
		for i := 0; i < len(instrNode.Content); i += 2 {
			key := instrNode.Content[i]
			val := instrNode.Content[i+1]
			switch key.Value {
			case "type":
				instrType = val.Value
			case "args":
				argsNode = val
			}
		}
		if instrType == "" {
			return Task{}, fmt.Errorf("instruction missing type")
		}
		if argsNode == nil {
			return Task{}, fmt.Errorf("instruction missing args")
		}

		var decodedArgs interface{}
		switch instrType {
		case "create_worker":
			var args CreateWorkerInstructionArgs
			if err := argsNode.Decode(&args); err != nil {
				return Task{}, fmt.Errorf("failed to decode create_worker args: %w", err)
			}
			decodedArgs = args

		case "remove_worker":
			var args RemoveWorkerInstructionArgs
			if err := argsNode.Decode(&args); err != nil {
				return Task{}, fmt.Errorf("failed to decode remove_worker args: %w", err)
			}
			decodedArgs = args

		case "stop_worker":
			var args StopWorkerInstructionArgs
			if err := argsNode.Decode(&args); err != nil {
				return Task{}, fmt.Errorf("failed to decode stop_worker args: %w", err)
			}
			decodedArgs = args

		case "start_worker":
			// Use a temporary struct to capture the raw YAML node.
			type tempStartArgs struct {
				WorkerUUID      uuid.UUID `yaml:"worker_uuid"`
				WorkerRawConfig yaml.Node `yaml:"worker_raw_config"`
			}
			var tempArgs tempStartArgs
			if err := argsNode.Decode(&tempArgs); err != nil {
				return Task{}, fmt.Errorf("failed to decode start_worker args: %w", err)
			}
			// Re-marshal the worker_raw_config node back into YAML bytes.
			rawConfigBytes, err := yaml.Marshal(&tempArgs.WorkerRawConfig)
			if err != nil {
				return Task{}, fmt.Errorf("failed to marshal worker_raw_config: %w", err)
			}
			startArgs := StartWorkerInstructionArgs{
				WorkerUUID:      tempArgs.WorkerUUID,
				WorkerRawConfig: rawConfigBytes,
			}
			decodedArgs = startArgs

		default:
			// For any unknown instruction type, decode the args into a generic map.
			var genericArgs map[string]interface{}
			if err := argsNode.Decode(&genericArgs); err != nil {
				return Task{}, fmt.Errorf("failed to decode generic args for instruction %s: %w", instrType, err)
			}
			decodedArgs = genericArgs
		}

		task.Instructions = append(task.Instructions, Instruction{
			Type: instrType,
			Args: decodedArgs,
		})
	}

	return task, nil
}

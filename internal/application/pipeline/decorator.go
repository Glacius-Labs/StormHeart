package pipeline

type Decorator func(TargetFunc) TargetFunc

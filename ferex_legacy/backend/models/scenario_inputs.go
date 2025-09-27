package models

// This file is intentionally left blank.
// The composite input structs (CoreScenarioInputs, FersScenarioInputs, etc.) previously defined here
// have been removed as their fields were directly incorporated into the ScenarioVariant struct
// in backend/scenario/file_format.go to simplify the data model for file saving/loading.
// This file will be populated with more higher-level composite input structs
// that are directly part of the ScenarioVariant's main input groups,
// following the 'fully independent variants (after creation)' model.
// Basic, shared types like ServicePeriod, LWOPPeriod, etc., belong in common_types.go.

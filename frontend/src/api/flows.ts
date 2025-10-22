import { apiClient } from "./client";

export interface FlowDefinition {
  nodes: Array<Record<string, unknown>>;
  edges: Array<Record<string, unknown>>;
}

export interface Flow {
  id: string;
  name: string;
  description: string;
  definition: FlowDefinition;
  metadata: Record<string, string>;
  version: number;
  updatedAt: string;
}

export interface CreateFlowInput {
  name: string;
  description: string;
  definition: FlowDefinition;
  metadata: Record<string, string>;
}

export const listFlows = async (): Promise<Flow[]> => {
  const response = await apiClient.get<Flow[]>("/flows");
  return response.data;
};

export const createFlow = async (payload: CreateFlowInput): Promise<Flow> => {
  const response = await apiClient.post<Flow>("/flows", payload);
  return response.data;
};

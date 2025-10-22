import { apiClient } from "./client";

export type WorkOrderStatus = "pending" | "running" | "failed" | "complete";

export interface WorkOrder {
  id: string;
  flowId: string;
  title: string;
  assignee: string;
  status: WorkOrderStatus;
  payload: Record<string, unknown>;
  metadata: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWorkOrderInput {
  flowId: string;
  title: string;
  assignee?: string;
  payload?: Record<string, unknown>;
  metadata?: Record<string, string>;
}

export const listWorkOrders = async (): Promise<WorkOrder[]> => {
  const response = await apiClient.get<WorkOrder[]>("/workorders");
  return response.data;
};

export const createWorkOrder = async (payload: CreateWorkOrderInput): Promise<WorkOrder> => {
  const response = await apiClient.post<WorkOrder>("/workorders", payload);
  return response.data;
};

export const retryWorkOrder = async (id: string): Promise<void> => {
  await apiClient.post(`/workorders/${id}/retry`);
};

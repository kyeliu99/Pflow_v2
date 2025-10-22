import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  CreateWorkOrderInput,
  WorkOrder,
  createWorkOrder,
  listWorkOrders,
  retryWorkOrder
} from "../api/workorders";

const workOrdersKey = ["workorders"];

export const useWorkOrders = () =>
  useQuery({
    queryKey: workOrdersKey,
    queryFn: listWorkOrders,
    refetchInterval: 5000
  });

export const useCreateWorkOrder = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateWorkOrderInput) => createWorkOrder(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workOrdersKey });
    }
  });
};

export const useRetryWorkOrder = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => retryWorkOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: workOrdersKey });
    }
  });
};

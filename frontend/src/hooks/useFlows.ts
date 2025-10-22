import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CreateFlowInput, Flow, createFlow, listFlows } from "../api/flows";

const flowsKey = ["flows"];

export const useFlows = () =>
  useQuery({
    queryKey: flowsKey,
    queryFn: listFlows,
    staleTime: 1000 * 30
  });

export const useCreateFlow = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateFlowInput) => createFlow(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: flowsKey });
    }
  });
};

export const useFlowOptions = () => {
  const { data } = useFlows();
  return (data ?? []).map((flow: Flow) => ({
    label: `${flow.name} v${flow.version}`,
    value: flow.id
  }));
};

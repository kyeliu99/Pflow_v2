import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  Chip,
  MenuItem,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useMemo, useState } from "react";
import { useFlowOptions } from "../hooks/useFlows";
import { useCreateWorkOrder, useRetryWorkOrder, useWorkOrders } from "../hooks/useWorkOrders";

const statusLabels: Record<string, string> = {
  pending: "待执行",
  running: "执行中",
  failed: "失败",
  complete: "完成"
};

function WorkOrderDashboard() {
  const flowOptions = useFlowOptions();
  const { data: workOrders } = useWorkOrders();
  const createWorkOrderMutation = useCreateWorkOrder();
  const retryMutation = useRetryWorkOrder();

  const [form, setForm] = useState({
    flowId: "",
    title: "自动化变更",
    assignee: "",
    payload: "{\n  \"context\": \"demo\"\n}"
  });

  const latest = useMemo(() => (workOrders ?? []).slice(0, 5), [workOrders]);

  const handleChange = (field: string) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, [field]: event.target.value }));
  };

  const handleCreate = () => {
    try {
      const payload = form.payload ? JSON.parse(form.payload) : {};
      createWorkOrderMutation.mutate({
        flowId: form.flowId,
        title: form.title,
        assignee: form.assignee,
        payload
      });
    } catch (error) {
      console.error("Invalid payload", error);
    }
  };

  return (
    <Card sx={{ height: "100%" }}>
      <CardContent>
        <Stack spacing={2}>
          <Typography variant="h6" fontWeight={600}>
            工单调度
          </Typography>
          <TextField
            select
            label="关联流程"
            value={form.flowId}
            onChange={handleChange("flowId")}
            helperText="选择已发布的流程定义"
          >
            {flowOptions.map((option) => (
              <MenuItem key={option.value} value={option.value}>
                {option.label}
              </MenuItem>
            ))}
          </TextField>
          <TextField label="工单标题" value={form.title} onChange={handleChange("title")} fullWidth />
          <TextField label="执行人" value={form.assignee} onChange={handleChange("assignee")} fullWidth />
          <TextField
            label="上下文 Payload"
            value={form.payload}
            onChange={handleChange("payload")}
            multiline
            minRows={4}
          />
          <Box>
            <Typography variant="subtitle2" color="text.secondary" gutterBottom>
              最近工单
            </Typography>
            <Stack spacing={1}>
              {latest.map((item) => (
                <Stack
                  direction="row"
                  alignItems="center"
                  justifyContent="space-between"
                  key={item.id}
                  sx={{
                    px: 1.5,
                    py: 1,
                    borderRadius: 1,
                    border: (theme) => `1px solid ${theme.palette.divider}`
                  }}
                >
                  <Box>
                    <Typography fontWeight={500}>{item.title}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {new Date(item.updatedAt).toLocaleString()} · {item.assignee || "未分配"}
                    </Typography>
                  </Box>
                  <Stack direction="row" spacing={1}>
                    <Chip label={statusLabels[item.status]} color={item.status === "failed" ? "error" : "primary"} />
                    {item.status === "failed" && (
                      <Button size="small" onClick={() => retryMutation.mutate(item.id)}>
                        重试
                      </Button>
                    )}
                  </Stack>
                </Stack>
              ))}
            </Stack>
          </Box>
        </Stack>
      </CardContent>
      <CardActions>
        <Button
          variant="contained"
          onClick={handleCreate}
          disabled={createWorkOrderMutation.isPending || !form.flowId}
        >
          {createWorkOrderMutation.isPending ? "创建中..." : "创建工单"}
        </Button>
      </CardActions>
    </Card>
  );
}

export default WorkOrderDashboard;

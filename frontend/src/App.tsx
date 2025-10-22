import { Box, Container, Grid, Stack, Typography } from "@mui/material";
import FlowDesigner from "./components/FlowDesigner";
import WorkOrderDashboard from "./components/WorkOrderDashboard";

function App() {
  return (
    <Container maxWidth={false} sx={{ py: 4 }}>
      <Stack spacing={3}>
        <Box>
          <Typography variant="h4" fontWeight={600} gutterBottom>
            PFlow 流程引擎 & 工单管理平台
          </Typography>
          <Typography variant="body1" color="text.secondary">
            拖拽式流程建模、工单生命周期管理与 Camunda BPMN2.0 引擎的深度集成。
          </Typography>
        </Box>
        <Grid container spacing={3}>
          <Grid item xs={12} md={7}>
            <FlowDesigner />
          </Grid>
          <Grid item xs={12} md={5}>
            <WorkOrderDashboard />
          </Grid>
        </Grid>
      </Stack>
    </Container>
  );
}

export default App;

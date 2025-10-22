import { useCallback, useMemo, useState } from "react";
import ReactFlow, {
  Background,
  Connection,
  Controls,
  Edge,
  MiniMap,
  Node,
  addEdge,
  useEdgesState,
  useNodesState
} from "react-flow-renderer";
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useCreateFlow, useFlows } from "../hooks/useFlows";

const initialNodes: Node[] = [
  {
    id: "start",
    type: "input",
    data: { label: "开始" },
    position: { x: 100, y: 80 }
  }
];

const initialEdges: Edge[] = [];

function FlowDesigner() {
  const { data: flows } = useFlows();
  const createFlowMutation = useCreateFlow();

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [name, setName] = useState("新流程");

  const definition = useMemo(
    () => ({
      nodes: nodes.map((node) => ({ id: node.id, type: node.type, position: node.position, data: node.data })),
      edges: edges.map((edge) => ({ id: edge.id, source: edge.source, target: edge.target }))
    }),
    [nodes, edges]
  );

  const onConnect = useCallback(
    (connection: Connection) => setEdges((eds) => addEdge(connection, eds)),
    [setEdges]
  );

  const handleSubmit = () => {
    createFlowMutation.mutate({
      name,
      description: "通过可视化设计器创建",
      definition,
      metadata: {
        designer: "react-flow",
        lastUpdated: new Date().toISOString()
      }
    });
  };

  return (
    <Card sx={{ height: "100%" }}>
      <CardContent sx={{ height: "100%" }}>
        <Stack spacing={2} height="100%">
          <Typography variant="h6" fontWeight={600}>
            流程建模
          </Typography>
          <TextField label="流程名称" value={name} onChange={(event) => setName(event.target.value)} fullWidth />
          <Box sx={{ height: 420, borderRadius: 2, overflow: "hidden", border: (theme) => `1px solid ${theme.palette.divider}` }}>
            <ReactFlow
              nodes={nodes}
              edges={edges}
              onNodesChange={onNodesChange}
              onEdgesChange={onEdgesChange}
              onConnect={onConnect}
              fitView
            >
              <MiniMap pannable zoomable />
              <Controls />
              <Background gap={16} size={1} />
            </ReactFlow>
          </Box>
          <Typography variant="body2" color="text.secondary">
            已发布流程：{flows?.length ?? 0}
          </Typography>
        </Stack>
      </CardContent>
      <CardActions>
        <Button variant="contained" onClick={handleSubmit} disabled={createFlowMutation.isPending}>
          {createFlowMutation.isPending ? "发布中..." : "发布流程"}
        </Button>
      </CardActions>
    </Card>
  );
}

export default FlowDesigner;

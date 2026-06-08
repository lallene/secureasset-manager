<script setup>
import { ref, onMounted } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

import { VueFlow } from "@vue-flow/core";
import { Background } from "@vue-flow/background";
import { Controls } from "@vue-flow/controls";
import { MiniMap } from "@vue-flow/minimap";

import "@vue-flow/core/dist/style.css";
import "@vue-flow/core/dist/theme-default.css";

const nodes = ref([]);
const edges = ref([]);
const error = ref("");

const selectedNode = ref(null);

const onNodeClick = (event) => {
  selectedNode.value = event.node.data;
};

const getToken = () => localStorage.getItem("token");

const fetchGraph = async () => {
  try {
    const response = await api.get("/cmdb/graph", {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    nodes.value = response.data.nodes || [];
    edges.value = response.data.edges || [];
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger le graphe CMDB";
  }
};

onMounted(fetchGraph);
</script>

<template>
  <DashboardLayout>
    <div class="graph-page">
      <div class="graph-header">
        <div>
          <p class="eyebrow">CMDB Mapping</p>
          <h1>CMDB Relationship Graph</h1>
          <p class="subtitle">
            Visualisation des dépendances entre services métiers, applications,
            bases de données et assets.
          </p>
        </div>
      </div>

      <div v-if="error" class="error-box">
        {{ error }}
      </div>

      <div class="graph-card">
        <VueFlow
            :nodes="nodes"
            :edges="edges"
            fit-view-on-init
            class="vue-flow-dark"
            @node-click="onNodeClick"
        >
          <Background />
          <Controls />
          <MiniMap />
        </VueFlow>
        <div v-if="selectedNode" class="detail-panel">
          <div class="panel-header">
            <p>{{ selectedNode.ci_type }}</p>
            <button @click="selectedNode = null">✕</button>
          </div>

          <h3>{{ selectedNode.name }}</h3>

          <div class="info-row" v-if="selectedNode.criticality">
            <span>Criticité</span>
            <strong>{{ selectedNode.criticality }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.status">
            <span>Statut</span>
            <strong>{{ selectedNode.status }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.version">
            <span>Version</span>
            <strong>{{ selectedNode.version }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.engine">
            <span>Moteur</span>
            <strong>{{ selectedNode.engine }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.environment">
            <span>Environnement</span>
            <strong>{{ selectedNode.environment }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.ip">
            <span>Adresse IP</span>
            <strong>{{ selectedNode.ip }}</strong>
          </div>

          <div class="info-row" v-if="selectedNode.site">
            <span>Site</span>
            <strong>{{ selectedNode.site }}</strong>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<style scoped>
.graph-page {
  padding: 2rem;
  min-height: 100vh;
  background: #0f172a;
  color: #f8fafc;
}

.graph-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 1.5rem;
}

.eyebrow {
  color: #818cf8;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.75rem;
  font-weight: 700;
  margin-bottom: 0.25rem;
}

h1 {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
}

.subtitle {
  color: #94a3b8;
  font-size: 0.9rem;
  margin-top: 0.35rem;
}

.graph-card {
  height: calc(100vh - 180px);
  background: #020617;
  border: 1px solid #334155;
  border-radius: 18px;
  overflow: hidden;
}

.error-box {
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid #ef4444;
  color: #fca5a5;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  margin-bottom: 1rem;
}

:deep(.vue-flow__node) {
  background: #1e293b;
  color: #f8fafc;
  border: 1px solid #475569;
  border-radius: 12px;
  padding: 10px 14px;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
}

:deep(.vue-flow__edge-path) {
  stroke: #818cf8;
  stroke-width: 2;
}

:deep(.vue-flow__controls) {
  background: #1e293b;
  border: 1px solid #334155;
}

:deep(.vue-flow__controls-button) {
  background: #1e293b;
  color: #f8fafc;
  border-bottom: 1px solid #334155;
}

:deep(.vue-flow__minimap) {
  background: #020617;
  border: 1px solid #334155;
}

.graph-card {
  position: relative;
}

.detail-panel {
  position: absolute;
  right: 18px;
  top: 18px;
  width: 320px;
  background: rgba(15, 23, 42, 0.96);
  border: 1px solid #334155;
  border-radius: 16px;
  padding: 18px;
  color: #f8fafc;
  z-index: 20;
  box-shadow: 0 20px 50px rgba(0,0,0,.45);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header p {
  margin: 0;
  font-size: 11px;
  color: #818cf8;
  text-transform: uppercase;
  letter-spacing: .08em;
  font-weight: 700;
}

.panel-header button {
  background: #1e293b;
  color: #94a3b8;
  border: 1px solid #334155;
  border-radius: 8px;
  cursor: pointer;
}

.detail-panel h3 {
  margin: 10px 0 16px;
  font-size: 18px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  border-top: 1px solid #334155;
  padding: 10px 0;
  font-size: 13px;
}

.info-row span {
  color: #94a3b8;
}

.info-row strong {
  color: #f8fafc;
  text-align: right;
}
</style>
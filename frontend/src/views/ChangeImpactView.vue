<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const impact = ref(null);
const error = ref("");

const getToken = () => localStorage.getItem("token");

const fetchImpact = async () => {
  try {
    loading.value = true;
    error.value = "";

    const response = await api.get(`/changes/${route.params.id}/impact`, {
      headers: {
        Authorization: `Bearer ${getToken()}`,
      },
    });

    impact.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger l'analyse d'impact";
  } finally {
    loading.value = false;
  }
};

const riskClass = (risk) => {
  return {
    Low: "risk-low",
    Medium: "risk-medium",
    High: "risk-high",
    Critical: "risk-critical",
  }[risk] || "risk-medium";
};

const statusClass = (status) => {
  return {
    Draft: "status-draft",
    Submitted: "status-submitted",
    Approved: "status-approved",
    Implemented: "status-implemented",
    Closed: "status-closed",
    Rejected: "status-rejected",
  }[status] || "status-draft";
};

const formatDate = (value) => {
  if (!value) return "Non défini";

  return new Date(value).toLocaleString("fr-FR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

onMounted(fetchImpact);
</script>

<template>
  <DashboardLayout>
    <div class="impact-page">
      <div class="page-top">
        <button class="back-btn" @click="router.push('/changes')">
          ← Retour aux changements
        </button>
      </div>

      <div v-if="loading" class="state-card">
        Chargement de l'analyse d'impact...
      </div>

      <div v-else-if="error" class="state-card error">
        {{ error }}
      </div>

      <template v-else-if="impact">
        <div class="header-card">
          <div class="header-main">
            <div class="ref">
              {{ impact.change.reference }}
            </div>

            <h1>
              {{ impact.change.title }}
            </h1>

            <p>
              {{ impact.change.description || "Aucune description renseignée." }}
            </p>
          </div>

          <div class="badges">
            <span class="badge" :class="riskClass(impact.change.risk)">
              Risque {{ impact.change.risk }}
            </span>

            <span class="badge" :class="statusClass(impact.change.status)">
              {{ impact.change.status }}
            </span>
          </div>
        </div>

        <div class="stats-grid">
          <div class="stat-card">
            <span>{{ impact.applications?.length || 0 }}</span>
            <small>Applications</small>
          </div>

          <div class="stat-card">
            <span>{{ impact.databases?.length || 0 }}</span>
            <small>Databases</small>
          </div>

          <div class="stat-card">
            <span>{{ impact.assets?.length || 0 }}</span>
            <small>Assets</small>
          </div>

          <div class="stat-card risk-stat">
            <span>{{ impact.change.risk }}</span>
            <small>Risk Level</small>
          </div>
        </div>

        <div class="impact-flow">
          <div class="flow-node business">
            <div class="flow-icon">🏢</div>
            <strong>{{ impact.business_service?.name || "Aucun service" }}</strong>
            <small>Business Service</small>
          </div>

          <div class="flow-arrow">→</div>

          <div class="flow-node application">
            <div class="flow-icon">📦</div>
            <strong>{{ impact.applications?.length || 0 }}</strong>
            <small>Applications</small>
          </div>

          <div class="flow-arrow">→</div>

          <div class="flow-node database">
            <div class="flow-icon">🗄</div>
            <strong>{{ impact.databases?.length || 0 }}</strong>
            <small>Databases</small>
          </div>

          <div class="flow-arrow">→</div>

          <div class="flow-node asset danger">
            <div class="flow-icon">💻</div>
            <strong>{{ impact.assets?.length || 0 }}</strong>
            <small>Assets</small>
          </div>
        </div>

        <div class="grid">
          <div class="card">
            <div class="card-header">
              <h2>🏢 Service Métier</h2>
            </div>

            <div v-if="impact.business_service" class="item">
              <div>
                <strong>{{ impact.business_service.name }}</strong>
                <small>{{ impact.business_service.description }}</small>
              </div>

              <span class="pill criticality">
                {{ impact.business_service.criticality }}
              </span>
            </div>

            <div v-else class="empty">
              Aucun service métier associé.
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              <h2>📦 Applications impactées</h2>
              <span>{{ impact.applications?.length || 0 }}</span>
            </div>

            <div
                v-for="app in impact.applications"
                :key="app.ID"
                class="item"
            >
              <div>
                <strong>{{ app.name }}</strong>
                <small>{{ app.description }}</small>
              </div>

              <span class="pill">
                v{{ app.version }}
              </span>
            </div>

            <div v-if="!impact.applications?.length" class="empty">
              Aucune application impactée.
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              <h2>🗄 Databases impactées</h2>
              <span>{{ impact.databases?.length || 0 }}</span>
            </div>

            <div
                v-for="db in impact.databases"
                :key="db.ID"
                class="item"
            >
              <div>
                <strong>{{ db.name }}</strong>
                <small>{{ db.environment }} · {{ db.site?.name || "Site inconnu" }}</small>
              </div>

              <span class="pill">
                {{ db.engine }} {{ db.version }}
              </span>
            </div>

            <div v-if="!impact.databases?.length" class="empty">
              Aucune base de données impactée.
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              <h2>💻 Assets impactés</h2>
              <span>{{ impact.assets?.length || 0 }}</span>
            </div>

            <div
                v-for="asset in impact.assets"
                :key="asset.ID"
                class="item"
            >
              <div>
                <strong>{{ asset.name }}</strong>
                <small>{{ asset.ip_address }} · {{ asset.site?.name || "Site inconnu" }}</small>
              </div>

              <span class="pill">
                {{ asset.type }}
              </span>
            </div>

            <div v-if="!impact.assets?.length" class="empty">
              Aucun asset impacté.
            </div>
          </div>
        </div>

        <div class="bottom-grid">
          <div class="card rollback-card">
            <h2>🔄 Rollback Plan</h2>

            <div class="rollback-content">
              {{ impact.change.rollback_plan || "Aucun plan de rollback renseigné." }}
            </div>
          </div>

          <div class="card timeline-card">
            <h2>📅 Timeline CAB</h2>

            <div class="timeline-row">
              <span>Créé le</span>
              <strong>{{ formatDate(impact.change.CreatedAt) }}</strong>
            </div>

            <div class="timeline-row">
              <span>Début planifié</span>
              <strong>{{ formatDate(impact.change.planned_start) }}</strong>
            </div>

            <div class="timeline-row">
              <span>Fin planifiée</span>
              <strong>{{ formatDate(impact.change.planned_end) }}</strong>
            </div>

            <div class="timeline-row">
              <span>Statut</span>
              <strong>{{ impact.change.status }}</strong>
            </div>
          </div>
        </div>
      </template>
    </div>
  </DashboardLayout>
</template>

<style scoped>
.impact-page {
  padding: 24px;
  background: #0f172a;
  min-height: 100vh;
  color: white;
}

.page-top {
  margin-bottom: 16px;
}

.back-btn {
  background: #1e293b;
  border: 1px solid #334155;
  color: #cbd5e1;
  padding: 10px 14px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
}

.back-btn:hover {
  background: #334155;
  color: #fff;
}

.state-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 16px;
  padding: 24px;
  color: #cbd5e1;
}

.state-card.error {
  border-color: #ef4444;
  color: #fca5a5;
  background: rgba(239, 68, 68, 0.1);
}

.header-card {
  background: linear-gradient(135deg, #1e293b, #111827);
  border: 1px solid #334155;
  border-radius: 18px;
  padding: 26px;
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  gap: 24px;
}

.header-main {
  min-width: 0;
}

.ref {
  color: #818cf8;
  font-weight: 800;
  margin-bottom: 8px;
  letter-spacing: 0.04em;
}

.header-card h1 {
  margin: 0;
  font-size: 28px;
  line-height: 1.2;
}

.header-card p {
  color: #94a3b8;
  margin-top: 10px;
  max-width: 900px;
}

.badges {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex-shrink: 0;
}

.badge {
  padding: 8px 14px;
  border-radius: 999px;
  text-align: center;
  font-size: 12px;
  font-weight: 800;
  white-space: nowrap;
}

.risk-critical {
  background: rgba(239, 68, 68, 0.18);
  color: #fecaca;
  border: 1px solid rgba(239, 68, 68, 0.35);
}

.risk-high {
  background: rgba(249, 115, 22, 0.18);
  color: #fdba74;
  border: 1px solid rgba(249, 115, 22, 0.35);
}

.risk-medium {
  background: rgba(234, 179, 8, 0.18);
  color: #fde68a;
  border: 1px solid rgba(234, 179, 8, 0.35);
}

.risk-low {
  background: rgba(34, 197, 94, 0.18);
  color: #bbf7d0;
  border: 1px solid rgba(34, 197, 94, 0.35);
}

.status-draft,
.status-closed {
  background: rgba(148, 163, 184, 0.18);
  color: #cbd5e1;
  border: 1px solid rgba(148, 163, 184, 0.35);
}

.status-submitted {
  background: rgba(59, 130, 246, 0.18);
  color: #93c5fd;
  border: 1px solid rgba(59, 130, 246, 0.35);
}

.status-approved {
  background: rgba(34, 197, 94, 0.18);
  color: #86efac;
  border: 1px solid rgba(34, 197, 94, 0.35);
}

.status-implemented {
  background: rgba(168, 85, 247, 0.18);
  color: #d8b4fe;
  border: 1px solid rgba(168, 85, 247, 0.35);
}

.status-rejected {
  background: rgba(239, 68, 68, 0.18);
  color: #fca5a5;
  border: 1px solid rgba(239, 68, 68, 0.35);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
}

.stat-card span {
  display: block;
  font-size: 28px;
  font-weight: 800;
}

.stat-card small {
  color: #94a3b8;
}

.risk-stat span {
  color: #fdba74;
  font-size: 22px;
}

.impact-flow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 18px;
  margin-bottom: 28px;
  background: #020617;
  border: 1px solid #334155;
  border-radius: 18px;
  padding: 18px;
  overflow-x: auto;
}

.flow-node {
  min-width: 150px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 14px;
  padding: 16px;
  text-align: center;
}

.flow-node strong {
  display: block;
  margin-top: 8px;
  font-size: 14px;
}

.flow-node small {
  color: #94a3b8;
}

.flow-icon {
  font-size: 24px;
}

.flow-arrow {
  font-size: 24px;
  color: #64748b;
}

.flow-node.business {
  border-color: rgba(168, 85, 247, 0.5);
}

.flow-node.application {
  border-color: rgba(59, 130, 246, 0.5);
}

.flow-node.database {
  border-color: rgba(34, 197, 94, 0.5);
}

.flow-node.asset {
  border-color: rgba(249, 115, 22, 0.5);
}

.flow-node.danger {
  box-shadow: 0 0 0 1px rgba(249, 115, 22, 0.2);
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 20px;
  margin-top: 20px;
}

.card {
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 16px;
  padding: 20px;
}

.card h2 {
  margin-top: 0;
  margin-bottom: 16px;
  font-size: 18px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  color: #94a3b8;
  font-size: 13px;
}

.item {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 12px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.item strong {
  color: white;
  display: block;
}

.item small {
  color: #64748b;
  display: block;
  margin-top: 4px;
  max-width: 560px;
}

.pill {
  color: #cbd5e1;
  background: #111827;
  border: 1px solid #334155;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 12px;
  white-space: nowrap;
}

.criticality {
  color: #fecaca;
  border-color: rgba(239, 68, 68, 0.4);
}

.empty {
  color: #94a3b8;
  padding: 12px;
  background: #0f172a;
  border: 1px dashed #334155;
  border-radius: 12px;
}

.rollback-content {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 12px;
  padding: 16px;
  color: #cbd5e1;
  line-height: 1.6;
}

.timeline-row {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  padding: 12px 0;
  border-bottom: 1px solid #334155;
}

.timeline-row:last-child {
  border-bottom: none;
}

.timeline-row span {
  color: #94a3b8;
}

.timeline-row strong {
  text-align: right;
}

@media (max-width: 1100px) {
  .stats-grid,
  .grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .header-card {
    flex-direction: column;
  }

  .impact-flow {
    justify-content: flex-start;
  }
}
</style>
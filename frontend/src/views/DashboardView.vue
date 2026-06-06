<script setup>
import { computed, onMounted, ref } from "vue";
import {
  Chart as ChartJS,
  ArcElement,
  BarElement,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler,
} from "chart.js";
import { Doughnut, Bar, Line } from "vue-chartjs";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

// ─── Chart.js registration ────────────────────────────────────────────────────
ChartJS.register(
  ArcElement,
  BarElement,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler
);

// ─── Constants ────────────────────────────────────────────────────────────────
const CATEGORY_COLORS = ["#534AB7", "#378ADD", "#1D9E75", "#EF9F27", "#888780"];

const CHART_STYLE = {
  color: "#888780",
  grid: "rgba(128,128,120,0.08)",
  font: { size: 10 },
};


const sites = ref([]);
const services = ref([]);

const filters = ref({
  site_id: "",
  service_id: "",
  period: "30",
});

const fetchSites = async () => {
  const token = localStorage.getItem("token");

  const response = await api.get("/sites", {
    headers: { Authorization: `Bearer ${token}` },
  });

  sites.value = response.data;
};

const fetchServices = async () => {
  const token = localStorage.getItem("token");

  const response = await api.get("/services", {
    headers: { Authorization: `Bearer ${token}` },
  });

  services.value = response.data;
};

// ─── State ────────────────────────────────────────────────────────────────────
const stats = ref({
  // Infrastructure
  total_assets: 0,
  active_assets: 0,
  alert_assets: 0,
  offline_assets: 0,
  // Incidents — statuts
  total_incidents: 0,
  open_incidents: 0,
  in_progress_incidents: 0,
  resolved_incidents: 0,
  closed_incidents: 0,
  // Incidents — priorités
  low_incidents: 0,
  medium_incidents: 0,
  high_incidents: 0,
  critical_incidents: 0,
  // MTTR par priorité (heures)
  mttr_critical: 0,
  mttr_high: 0,
  mttr_medium: 0,
  mttr_low: 0,
  // SLA & qualité
  overdue_incidents: 0,
  sla_resolution_rate: 0,
  mttr_hours: 0,
  reopening_rate: 0,
  // Séries temporelles
  weekly_trend_opened: [],
  weekly_trend_resolved: [],
  monthly_created: [],
  monthly_resolved: [],
  sla_history: [],
  // Relations
  incident_categories: [],  // [{ type: string, count: number }]
  service_workload: [],     // [{ service: string, open: number, in_progress: number }]
  // Divers
  unread_notifications: 0,
});

// ─── API ──────────────────────────────────────────────────────────────────────
const fetchStats = async () => {
  try {
    const token = localStorage.getItem("token");

    const response = await api.get("/dashboard/stats", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      params: {
        site_id: filters.value.site_id || undefined,
        service_id: filters.value.service_id || undefined,
        period: filters.value.period,
      },
    });

    stats.value = response.data;

    console.log("Dashboard Stats :", response.data);
  } catch (error) {
    console.error("[Dashboard] Erreur fetchStats :", error);
  }
};

// ─── Computed — Infrastructure ────────────────────────────────────────────────
const assetPercent = (key) =>
  computed(() =>
    stats.value.total_assets
      ? Math.round((stats.value[key] / stats.value.total_assets) * 100)
      : 0
  );

const activeAssetsPercent  = assetPercent("active_assets");
const alertAssetsPercent   = assetPercent("alert_assets");
const offlineAssetsPercent = assetPercent("offline_assets");

// ─── Helpers ──────────────────────────────────────────────────────────────────
const fmt = (value, decimals = 2) =>
  value != null ? Number(value).toFixed(decimals) : "0";

// ─── Chart base options ───────────────────────────────────────────────────────
const baseOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  layout: { padding: { top: 4, bottom: 4 } },
};

const tickStyle = { ...CHART_STYLE.font, color: CHART_STYLE.color };

const commonScales = {
  x: { grid: { display: false }, ticks: tickStyle },
  y: { grid: { color: CHART_STYLE.grid }, ticks: tickStyle },
};

const stackedScales = {
  x: { stacked: true, grid: { display: false }, ticks: tickStyle },
  y: { stacked: true, grid: { color: CHART_STYLE.grid }, ticks: tickStyle },
};

// ─── Chart data (computed) ────────────────────────────────────────────────────

/** C1 — Tendance hebdomadaire */
const c1Data = computed(() => ({
  labels: ["Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"],
  datasets: [
    {
      label: "Ouverts",
      data: stats.value.weekly_trend_opened,
      borderColor: "#534AB7",
      backgroundColor: "rgba(83,74,183,.08)",
      fill: true,
      tension: 0.4,
      pointRadius: 3,
      pointBackgroundColor: "#534AB7",
      borderWidth: 2,
    },
    {
      label: "Résolus",
      data: stats.value.weekly_trend_resolved,
      borderColor: "#1D9E75",
      backgroundColor: "rgba(29,158,117,.06)",
      fill: true,
      tension: 0.4,
      pointRadius: 3,
      pointBackgroundColor: "#1D9E75",
      borderWidth: 2,
    },
  ],
}));

/** C2 — Répartition par statut (donut) */
const c2Data = computed(() => ({
  labels: ["Ouverts", "En cours", "Résolus", "Clôturés"],
  datasets: [
    {
      data: [
        stats.value.open_incidents,
        stats.value.in_progress_incidents,
        stats.value.resolved_incidents,
        stats.value.closed_incidents,
      ],
      backgroundColor: ["#D85A30", "#534AB7", "#1D9E75", "#888780"],
      borderWidth: 0,
      hoverOffset: 4,
    },
  ],
}));

/** C3 — Volume mensuel créés vs résolus */
const c3Data = computed(() => ({
  labels: ["Jan", "Fév", "Mar", "Avr", "Mai", "Jun"],
  datasets: [
    {
      label: "Créés",
      data: stats.value.monthly_created,
      backgroundColor: "#378ADD",
      borderRadius: 3,
      borderSkipped: false,
    },
    {
      label: "Résolus",
      data: stats.value.monthly_resolved,
      backgroundColor: "#1D9E75",
      borderRadius: 3,
      borderSkipped: false,
    },
  ],
}));

/** C4 — Priorité des incidents (barres horizontales) */
const c4Data = computed(() => ({
  labels: ["Faible", "Moyenne", "Haute", "Critique"],
  datasets: [
    {
      data: [
        stats.value.low_incidents,
        stats.value.medium_incidents,
        stats.value.high_incidents,
        stats.value.critical_incidents,
      ],
      backgroundColor: ["#888780", "#378ADD", "#EF9F27", "#E24B4A"],
      borderRadius: 4,
      borderSkipped: false,
    },
  ],
}));

const c4Options = {
  ...baseOptions,
  indexAxis: "y",
  scales: {
    x: { grid: { color: CHART_STYLE.grid }, ticks: tickStyle },
    y: { grid: { display: false }, ticks: tickStyle },
  },
};

/** C5 — MTTR par priorité */
const c5Data = computed(() => ({
  labels: ["Critique", "Haute", "Moyenne", "Faible"],
  datasets: [
    {
      data: [
        stats.value.mttr_critical,
        stats.value.mttr_high,
        stats.value.mttr_medium,
        stats.value.mttr_low,
      ],
      backgroundColor: ["#E24B4A", "#EF9F27", "#378ADD", "#888780"],
      borderRadius: 4,
      borderSkipped: false,
    },
  ],
}));

const c5Options = {
  ...baseOptions,
  scales: {
    ...commonScales,
    y: {
      grid: { color: CHART_STYLE.grid },
      ticks: { ...tickStyle, callback: (v) => v + "h" },
    },
  },
};

/** C6 — Incidents par catégorie (donut) */
const c6Data = computed(() => ({
  labels: stats.value.incident_categories?.map((i) => i.type) ?? [],
  datasets: [
    {
      data: stats.value.incident_categories?.map((i) => i.count) ?? [],
      backgroundColor: CATEGORY_COLORS,
      borderWidth: 0,
      hoverOffset: 3,
    },
  ],
}));

/** C7 — Charge par service (barres empilées) */
const c7Data = computed(() => ({
  labels: stats.value.service_workload?.map((i) => i.service) ?? [],
  datasets: [
    {
      label: "Ouverts",
      data: stats.value.service_workload?.map((i) => i.open) ?? [],
      backgroundColor: "#D85A30",
      borderRadius: 3,
      borderSkipped: false,
    },
    {
      label: "En cours",
      data: stats.value.service_workload?.map((i) => i.in_progress) ?? [],
      backgroundColor: "#534AB7",
      borderRadius: 3,
      borderSkipped: false,
    },
  ],
}));

/** C8 — Évolution SLA sur 12 semaines */
const c8Data = computed(() => ({
  labels: ["S12","S13","S14","S15","S16","S17","S18","S19","S20","S21","S22","S23"],
  datasets: [
    {
      label: "Taux SLA",
      data: stats.value.sla_history,
      borderColor: "#378ADD",
      backgroundColor: "rgba(55,138,221,.07)",
      fill: true,
      tension: 0.4,
      pointRadius: 3,
      pointBackgroundColor: "#378ADD",
      borderWidth: 2,
    },
    {
      label: "Objectif",
      data: Array(12).fill(95),
      borderColor: "#E24B4A",
      borderDash: [5, 4],
      borderWidth: 1.5,
      pointRadius: 0,
      fill: false,
    },
  ],
}));

const c8Options = {
  ...baseOptions,
  scales: {
    x: {
      grid: { display: false },
      ticks: { ...tickStyle, autoSkip: false, maxRotation: 0 },
    },
    y: {
      min: 0,
      max: 100,
      grid: { color: CHART_STYLE.grid },
      ticks: { ...tickStyle, callback: (v) => v + "%" },
    },
  },
};

const resolvedIncidentPercent = computed(() =>
    stats.value.total_incidents
        ? Math.round(
            ((stats.value.resolved_incidents + stats.value.closed_incidents) /
                stats.value.total_incidents) *
            100
        )
        : 0
);

// ─── Lifecycle ────────────────────────────────────────────────────────────────
onMounted(() => {
  fetchSites();
  fetchServices();
  fetchStats();
});
</script>

<template>
  <DashboardLayout>
    <div class="db antialiased">
      <h2 class="sr-only">Tableau de bord ITSM — métriques, tendances et analyse des incidents</h2>

      <!-- ── Topbar ──────────────────────────────────────────────────────────── -->
      <div class="topbar">
        <div>
          <h1>
            <i class="ti ti-layout-dashboard" aria-hidden="true" />
            Tableau de bord opérationnel ITSM
          </h1>
          <p>Infrastructure &amp; incidents · Semaine en cours · Mis à jour aujourd'hui</p>
        </div>
        <div class="filters">
          <select v-model="filters.site_id" @change="fetchStats">
            <option value="">Tous les sites</option>
            <option
                v-for="site in sites"
                :key="site.ID"
                :value="site.ID"
            >
              {{ site.name }}
            </option>
          </select>

          <select v-model="filters.service_id" @change="fetchStats">
            <option value="">Tous les services</option>
            <option
                v-for="service in services"
                :key="service.ID"
                :value="service.ID"
            >
              {{ service.name }}
            </option>
          </select>

          <select v-model="filters.period" @change="fetchStats">
            <option value="7">7 derniers jours</option>
            <option value="30">30 derniers jours</option>
            <option value="90">90 derniers jours</option>
            <option value="180">6 derniers mois</option>
          </select>
        </div>

        <div class="badges">
          <span v-if="stats.critical_incidents > 0" class="badge b-red pulse">
            <i class="ti ti-alert-circle" aria-hidden="true" />
            {{ stats.critical_incidents }} critiques
          </span>
          <span v-if="stats.overdue_incidents > 0" class="badge b-amber">
            {{ stats.overdue_incidents }} hors SLA
          </span>
          <span v-if="stats.unread_notifications > 0" class="badge b-blue">
            {{ stats.unread_notifications }} notifications
          </span>
          <span class="badge b-green">{{ fmt(stats.sla_resolution_rate) }}% SLA</span>
        </div>
      </div>

      <!-- ── Section : Infrastructure ──────────────────────────────────────── -->
      <p class="sec">Infrastructure</p>

      <div class="krow">
        <div class="kc">
          <div class="kc-lbl">Parc total</div>
          <div class="kc-val">{{ stats.total_assets }}</div>
          <div class="kc-sub">
            {{ stats.active_assets }} actifs • {{ stats.alert_assets }} en alerte • {{ stats.offline_assets }} hors ligne
          </div>
          <div class="kc-bar">
            <div
                class="kc-fill"
                :style="{ width: activeAssetsPercent + '%', background: '#378ADD' }"
            />
          </div>
        </div>

        <div class="kc">
          <div class="kc-lbl">Assets actifs</div>
          <div class="kc-val">{{ stats.active_assets }}</div>
          <div class="kc-sub">{{ activeAssetsPercent }}% du parc disponible</div>
          <div class="kc-bar">
            <div
                class="kc-fill"
                :style="{ width: activeAssetsPercent + '%', background: '#1D9E75' }"
            />
          </div>
        </div>

        <div class="kc">
          <div class="kc-lbl">Assets en alerte</div>
          <div class="kc-val">{{ stats.alert_assets }}</div>
          <div class="kc-sub">{{ alertAssetsPercent }}% nécessitent une attention</div>
          <div class="kc-bar">
            <div
                class="kc-fill"
                :style="{ width: alertAssetsPercent + '%', background: '#EF9F27' }"
            />
          </div>
        </div>

        <div class="kc">
          <div class="kc-lbl">Assets hors ligne</div>
          <div class="kc-val">{{ stats.offline_assets }}</div>
          <div class="kc-sub">{{ offlineAssetsPercent }}% indisponibles</div>
          <div class="kc-bar">
            <div
                class="kc-fill"
                :style="{ width: offlineAssetsPercent + '%', background: '#E24B4A' }"
            />
          </div>
        </div>

        <div class="kc">
          <div class="kc-lbl">Incidents filtrés</div>
          <div class="kc-val">{{ stats.total_incidents }}</div>
          <div class="kc-sub">
            {{ stats.open_incidents }} ouverts • {{ stats.in_progress_incidents }} en cours
          </div>
          <div class="kc-bar">
            <div
                class="kc-fill"
                :style="{ width: resolvedIncidentPercent + '%', background: '#7F77DD' }"
            />
          </div>
        </div>
      </div>

      <!-- ── Section : Statuts ──────────────────────────────────────────────── -->
      <p class="sec">Statuts des incidents</p>

      <div class="srow">
        <div class="sc or">
          <div class="sc-l">Ouverts</div>
          <div class="sc-v">{{ stats.open_incidents }}</div>
          <span class="tag t-or">À traiter</span>
        </div>
        <div class="sc in">
          <div class="sc-l">En cours</div>
          <div class="sc-v">{{ stats.in_progress_incidents }}</div>
          <span class="tag t-in">Assignés</span>
        </div>
        <div class="sc em">
          <div class="sc-l">Résolus</div>
          <div class="sc-v">{{ stats.resolved_incidents }}</div>
          <span class="tag t-em">À valider</span>
        </div>
        <div class="sc sl">
          <div class="sc-l">Clôturés</div>
          <div class="sc-v">{{ stats.closed_incidents }}</div>
          <span class="tag t-sl">Archivés</span>
        </div>
        <div class="sc rs">
          <div class="sc-l">Hors SLA</div>
          <div class="sc-v">{{ stats.overdue_incidents }}</div>
          <span class="tag t-rs pulse">En retard</span>
        </div>
      </div>

      <!-- ── Section : Graphiques ───────────────────────────────────────────── -->
      <p class="sec">Analyse graphique</p>

      <!-- Ligne 1 : tendance + donut statuts -->
      <div class="g2">
        <div class="card">
          <h3>Tendance hebdomadaire — incidents ouverts</h3>
          <p>7 derniers jours · nouveaux tickets par jour</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#534AB7;" />Ouverts</span>
            <span class="li"><span class="ld" style="background:#1D9E75;" />Résolus</span>
          </div>
          <div class="chart-wrap">
            <Line :data="c1Data" :options="{ ...baseOptions, scales: commonScales }" />
          </div>
        </div>

        <div class="card">
          <h3>Répartition par statut</h3>
          <p>Distribution actuelle des {{ stats.total_incidents }} incidents</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#D85A30;" />Ouverts {{ stats.open_incidents }}</span>
            <span class="li"><span class="ld" style="background:#534AB7;" />En cours {{ stats.in_progress_incidents }}</span>
            <span class="li"><span class="ld" style="background:#1D9E75;" />Résolus {{ stats.resolved_incidents }}</span>
            <span class="li"><span class="ld" style="background:#888780;" />Clôturés {{ stats.closed_incidents }}</span>
          </div>
          <div class="chart-wrap">
            <Doughnut :data="c2Data" :options="{ ...baseOptions, cutout: '68%' }" />
          </div>
        </div>
      </div>

      <!-- Ligne 2 : volume mensuel + priorités -->
      <div class="g2">
        <div class="card">
          <h3>Volume mensuel — 6 derniers mois</h3>
          <p>Incidents créés vs résolus par mois</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#378ADD;" />Créés</span>
            <span class="li"><span class="ld" style="background:#1D9E75;" />Résolus</span>
          </div>
          <div class="chart-wrap">
            <Bar :data="c3Data" :options="{ ...baseOptions, scales: commonScales }" />
          </div>
        </div>

        <div class="card">
          <h3>Priorité des incidents actifs</h3>
          <p>Répartition des incidents non clôturés</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#888780;" />Faible {{ stats.low_incidents }}</span>
            <span class="li"><span class="ld" style="background:#378ADD;" />Moyenne {{ stats.medium_incidents }}</span>
            <span class="li"><span class="ld" style="background:#EF9F27;" />Haute {{ stats.high_incidents }}</span>
            <span class="li"><span class="ld" style="background:#E24B4A;" />Critique {{ stats.critical_incidents }}</span>
          </div>
          <div class="chart-wrap">
            <Bar :data="c4Data" :options="c4Options" />
          </div>
        </div>
      </div>

      <!-- Ligne 3 : MTTR + catégories + charge service -->
      <div class="g3">
        <div class="card">
          <h3>MTTR par priorité</h3>
          <p>Temps moyen de résolution (heures)</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#E24B4A;" />Critique</span>
            <span class="li"><span class="ld" style="background:#EF9F27;" />Haute</span>
            <span class="li"><span class="ld" style="background:#378ADD;" />Moyenne</span>
          </div>
          <div class="chart-wrap">
            <Bar :data="c5Data" :options="c5Options" />
          </div>
        </div>

        <div class="card">
          <h3>Incidents par catégorie</h3>
          <p>Nature des tickets ce mois</p>
          <div class="leg">
            <span
              v-for="(cat, i) in stats.incident_categories"
              :key="cat.type"
              class="li"
            >
              <span class="ld" :style="{ background: CATEGORY_COLORS[i % CATEGORY_COLORS.length] }" />
              {{ cat.type }} {{ cat.count }}
            </span>
          </div>
          <div class="chart-wrap">
            <Doughnut :data="c6Data" :options="{ ...baseOptions, cutout: '60%' }" />
          </div>
        </div>

        <div class="card">
          <h3>Charge par service</h3>
          <p>Tickets ouverts et en cours par service</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#D85A30;" />Ouverts</span>
            <span class="li"><span class="ld" style="background:#534AB7;" />En cours</span>
          </div>
          <div class="chart-wrap">
            <Bar :data="c7Data" :options="{ ...baseOptions, scales: stackedScales }" />
          </div>
        </div>
      </div>

      <!-- Ligne 4 : évolution SLA (pleine largeur) -->
      <div class="g1">
        <div class="card">
          <h3>Évolution SLA hebdomadaire — 12 dernières semaines</h3>
          <p>Taux de respect SLA (%) · Objectif 95%</p>
          <div class="leg">
            <span class="li"><span class="ld" style="background:#378ADD;" />Taux SLA réel</span>
            <span class="li">
              <span class="ld" style="background:#E24B4A; width:16px; height:2px; border-radius:0;" />
              Objectif 95%
            </span>
          </div>
          <div class="chart-wrap">
            <Line :data="c8Data" :options="c8Options" />
          </div>
        </div>
      </div>

      <!-- ── Section : Performance SLA ─────────────────────────────────────── -->
      <p class="sec">Performance SLA</p>

      <div class="sla-grid">
        <div class="sla-c">
          <div class="sla-l">Taux résolution SLA</div>
          <div class="sla-v">{{ fmt(stats.sla_resolution_rate) }}%</div>
          <div class="sla-bar">
            <div class="sla-fill sla-warn" :style="{ width: stats.sla_resolution_rate + '%' }" />
          </div>
          <div class="sla-d">Objectif 95%</div>
        </div>

        <div class="sla-c">
          <div class="sla-l">MTTR moyen</div>
          <div class="sla-v">{{ fmt(stats.mttr_hours, 1) }}h</div>
          <div class="sla-bar">
            <div class="sla-fill sla-warn" style="width:60%;" />
          </div>
          <div class="sla-d">Cible maximale : 3h</div>
        </div>

        <div class="sla-c">
          <div class="sla-l">Taux réouverture</div>
          <div class="sla-v">{{ fmt(stats.reopening_rate) }}%</div>
          <div class="sla-bar">
            <div
              class="sla-fill"
              :class="stats.reopening_rate < 5 ? 'sla-ok' : 'sla-bad'"
              :style="{ width: (100 - stats.reopening_rate) + '%' }"
            />
          </div>
          <div class="sla-d">Dans l'objectif &lt; 5%</div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>

<style scoped>
/* ── Reset & base ─────────────────────────────────────────────────────────── */
.db { padding: 0 0 1rem; }

.filters {
  display: flex;
  gap: 8px;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.filters select {
  background: var(--color-background-primary);
  border: 0.5px solid var(--color-border-tertiary);
  border-radius: var(--border-radius-md);
  padding: 8px 10px;
  font-size: 12px;
  color: var(--color-text-primary);
}

/* ── Topbar ───────────────────────────────────────────────────────────────── */
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.25rem 0 1rem;
  border-bottom: 0.5px solid var(--color-border-tertiary);
  margin-bottom: 1.25rem;
}
.topbar h1 {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 17px;
  font-weight: 500;
  color: var(--color-text-primary);
  letter-spacing: -0.3px;
}
.topbar p { font-size: 12px; color: var(--color-text-secondary); margin-top: 2px; }

/* ── Badges ───────────────────────────────────────────────────────────────── */
.badges { display: flex; gap: 6px; flex-wrap: wrap; align-items: center; }
.badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  font-weight: 500;
  padding: 3px 9px;
  border-radius: 100px;
  border: 0.5px solid;
}
.b-red    { background: #FCEBEB; color: #A32D2D; border-color: #F09595; }
.b-amber  { background: #FAEEDA; color: #854F0B; border-color: #FAC775; }
.b-blue   { background: #E6F1FB; color: #185FA5; border-color: #B5D4F4; }
.b-green  { background: #EAF3DE; color: #3B6D11; border-color: #C0DD97; }

/* ── Section title ────────────────────────────────────────────────────────── */
.sec {
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-tertiary);
  margin: 1.25rem 0 0.6rem;
}

/* ── KPI cards ────────────────────────────────────────────────────────────── */
.krow { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; margin-bottom: 1.25rem; }
.kc   { background: var(--color-background-secondary); border-radius: var(--border-radius-md); padding: 12px 14px; }
.kc-lbl { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 5px; }
.kc-val { font-size: 26px; font-weight: 500; color: var(--color-text-primary); letter-spacing: -1px; line-height: 1; }
.kc-sub { font-size: 11px; color: var(--color-text-tertiary); margin-top: 4px; }
.kc-bar { height: 2px; background: var(--color-border-tertiary); border-radius: 2px; margin-top: 8px; overflow: hidden; }
.kc-fill { height: 100%; border-radius: 2px; }

/* ── Status cards ─────────────────────────────────────────────────────────── */
.srow { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; margin-bottom: 1.25rem; }
.sc {
  background: var(--color-background-primary);
  border: 0.5px solid var(--color-border-tertiary);
  border-radius: var(--border-radius-md);
  padding: 11px 13px;
  position: relative;
  overflow: hidden;
}
.sc::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; }
.sc.or::before { background: #D85A30; }
.sc.in::before { background: #534AB7; }
.sc.em::before { background: #1D9E75; }
.sc.sl::before { background: #888780; }
.sc.rs::before { background: #A32D2D; }
.sc-l { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 6px; }
.sc-v { font-size: 26px; font-weight: 500; color: var(--color-text-primary); letter-spacing: -1px; line-height: 1; }
.tag  { display: inline-block; font-size: 10px; font-weight: 500; padding: 2px 7px; border-radius: 100px; margin-top: 6px; }
.t-or { background: #FAECE7; color: #993C1D; }
.t-in { background: #EEEDFE; color: #3C3489; }
.t-em { background: #EAF3DE; color: #3B6D11; }
.t-sl { background: #F1EFE8; color: #5F5E5A; }
.t-rs { background: #FCEBEB; color: #A32D2D; }

/* ── Chart grids ──────────────────────────────────────────────────────────── */
.g1 { display: grid; grid-template-columns: 1fr;       gap: 10px; margin-bottom: 10px; }
.g2 { display: grid; grid-template-columns: 1fr 1fr;   gap: 10px; margin-bottom: 10px; }
.g3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 10px; margin-bottom: 10px; }

/* ── Chart card ───────────────────────────────────────────────────────────── */
.card {
  background: var(--color-background-primary);
  border: 0.5px solid var(--color-border-tertiary);
  border-radius: var(--border-radius-lg);
  padding: 14px 16px;
}
.card h3 { font-size: 13px; font-weight: 500; color: var(--color-text-primary); margin-bottom: 2px; }
.card p  { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 10px; }
.chart-wrap { position: relative; height: 170px; }
.g3 .chart-wrap { height: 160px; }
.g1 .chart-wrap { height: 160px; }

/* ── Legend ───────────────────────────────────────────────────────────────── */
.leg { display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 8px; }
.li  { display: flex; align-items: center; gap: 5px; font-size: 11px; color: var(--color-text-secondary); }
.ld  { width: 8px; height: 8px; border-radius: 2px; flex-shrink: 0; }

/* ── SLA section ──────────────────────────────────────────────────────────── */
.sla-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin-bottom: 1.25rem; }
.sla-c    { background: var(--color-background-secondary); border-radius: var(--border-radius-md); padding: 12px 14px; }
.sla-l    { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 4px; }
.sla-v    { font-size: 20px; font-weight: 500; color: var(--color-text-primary); }
.sla-d    { font-size: 11px; color: var(--color-text-tertiary); margin-top: 3px; }
.sla-bar  { height: 4px; background: var(--color-border-tertiary); border-radius: 2px; margin-top: 7px; overflow: hidden; }
.sla-fill { height: 100%; border-radius: 2px; }
.sla-ok   { background: #1D9E75; }
.sla-warn { background: #EF9F27; }
.sla-bad  { background: #E24B4A; }

/* ── Animations ───────────────────────────────────────────────────────────── */
@keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
.pulse { animation: pulse 1.8s infinite; }
</style>
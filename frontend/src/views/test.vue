<script setup>
import { computed, onMounted, ref, watch } from "vue";
import {
  Chart as ChartJS, ArcElement, BarElement, CategoryScale, LinearScale,
  PointElement, LineElement, Tooltip, Legend, Filler,
} from "chart.js";
import { Doughnut, Bar, Line } from "vue-chartjs";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

ChartJS.register(
    ArcElement, BarElement, CategoryScale, LinearScale,
    PointElement, LineElement, Tooltip, Legend, Filler
);

// ─── Palette ──────────────────────────────────────────────────────────────────
const C = {
  indigo:  "#6366f1",
  cyan:    "#06b6d4",
  emerald: "#10b981",
  amber:   "#f59e0b",
  orange:  "#f97316",
  rose:    "#ef4444",
  slate:   "#94a3b8",
  violet:  "#8b5cf6",
};
const CAT_COLORS = [C.indigo, C.cyan, C.emerald, C.amber, C.slate];

// ─── Chart.js config commune ──────────────────────────────────────────────────
const TICK = { size: 10, family: "'DM Sans', sans-serif", color: "#94a3b8" };
const GRID = "rgba(148,163,184,0.07)";

const baseOpts = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  layout: { padding: { top: 4, bottom: 4 } },
};
const tickStyle = { ...TICK };
const commonScales = {
  x: { grid: { display: false }, ticks: tickStyle },
  y: { grid: { color: GRID },    ticks: tickStyle, border: { display: false } },
};
const stackedScales = {
  x: { stacked: true, grid: { display: false }, ticks: tickStyle },
  y: { stacked: true, grid: { color: GRID },    ticks: tickStyle, border: { display: false } },
};

// ─── Données filtres ──────────────────────────────────────────────────────────
const sites    = ref([]);
const services = ref([]);
const filters  = ref({ site_id: "", service_id: "", period: "30" });

const PERIOD_OPTS = [
  { label: "7 j",    value: "7"   },
  { label: "30 j",   value: "30"  },
  { label: "90 j",   value: "90"  },
  { label: "6 mois", value: "180" },
];

// ─── État principal ───────────────────────────────────────────────────────────
const isLoading   = ref(false);
const lastUpdated = ref(new Date());

const stats = ref({
  total_assets: 0, active_assets: 0, alert_assets: 0, offline_assets: 0,
  total_incidents: 0, open_incidents: 0, in_progress_incidents: 0,
  resolved_incidents: 0, closed_incidents: 0,
  low_incidents: 0, medium_incidents: 0, high_incidents: 0, critical_incidents: 0,
  mttr_critical: 0, mttr_high: 0, mttr_medium: 0, mttr_low: 0,
  overdue_incidents: 0, sla_resolution_rate: 0, mttr_hours: 0, reopening_rate: 0,
  unread_notifications: 0,
  weekly_trend_opened: [], weekly_trend_resolved: [],
  monthly_created: [], monthly_resolved: [],
  sla_history: [],
  incident_categories: [],
  service_workload: [],
  total_applications: 0, critical_applications: 0,
  total_databases: 0, production_databases: 0,
  total_relations: 0,
});

// ─── API ──────────────────────────────────────────────────────────────────────
const authHeader = () => ({ Authorization: `Bearer ${localStorage.getItem("token")}` });

const fetchSites = async () => {
  const { data } = await api.get("/sites", { headers: authHeader() });
  sites.value = data;
};
const fetchServices = async () => {
  const { data } = await api.get("/services", { headers: authHeader() });
  services.value = data;
};
const fetchStats = async () => {
  isLoading.value = true;
  try {
    const { data } = await api.get("/dashboard/stats", {
      headers: authHeader(),
      params: {
        site_id:    filters.value.site_id    || undefined,
        service_id: filters.value.service_id || undefined,
        period:     filters.value.period,
      },
    });
    stats.value   = data;
    lastUpdated.value = new Date();
  } catch (err) {
    console.error("[Dashboard] fetchStats:", err);
  } finally {
    isLoading.value = false;
  }
};

watch(filters, fetchStats, { deep: true });

// ─── Computed ─────────────────────────────────────────────────────────────────
const pct = (k) => computed(() =>
    stats.value.total_assets ? Math.round((stats.value[k] / stats.value.total_assets) * 100) : 0
);
const activeAssetsPercent  = pct("active_assets");
const alertAssetsPercent   = pct("alert_assets");
const offlineAssetsPercent = pct("offline_assets");

const resolvedPct = computed(() =>
    stats.value.total_incidents
        ? Math.round(((stats.value.resolved_incidents + stats.value.closed_incidents) / stats.value.total_incidents) * 100)
        : 0
);

const slaColor  = computed(() => stats.value.sla_resolution_rate >= 90 ? C.emerald : stats.value.sla_resolution_rate >= 75 ? C.amber : C.rose);
const mttrColor = computed(() => stats.value.mttr_hours <= 3 ? C.emerald : stats.value.mttr_hours <= 6 ? C.amber : C.rose);

// ─── Helpers ──────────────────────────────────────────────────────────────────
const fmt     = (v, d = 1) => v != null ? Number(v).toFixed(d) : "0";
const fmtTime = (h) => { if (!h) return "—"; const hrs = Math.floor(h); const m = Math.round((h - hrs) * 60); return m ? `${hrs}h ${m}m` : `${hrs}h`; };
const fmtClock = (d) => d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" });

// ─── Datasets Chart.js ────────────────────────────────────────────────────────
const c1Data = computed(() => ({
  labels: ["Lun","Mar","Mer","Jeu","Ven","Sam","Dim"],
  datasets: [
    { data: stats.value.weekly_trend_opened,   borderColor: C.indigo,  backgroundColor: "rgba(99,102,241,.07)",  fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.indigo,  borderWidth: 2 },
    { data: stats.value.weekly_trend_resolved, borderColor: C.emerald, backgroundColor: "rgba(16,185,129,.06)", fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.emerald, borderWidth: 2 },
  ],
}));

const c2Data = computed(() => ({
  labels: ["Ouverts","En cours","Résolus","Clôturés"],
  datasets: [{ data: [stats.value.open_incidents, stats.value.in_progress_incidents, stats.value.resolved_incidents, stats.value.closed_incidents], backgroundColor: [C.orange, C.indigo, C.emerald, C.slate], borderWidth: 0, hoverOffset: 6 }],
}));

const c3Data = computed(() => ({
  labels: ["Jan","Fév","Mar","Avr","Mai","Jun"],
  datasets: [
    { label: "Créés",   data: stats.value.monthly_created,  backgroundColor: C.indigo,  borderRadius: 4, borderSkipped: false },
    { label: "Résolus", data: stats.value.monthly_resolved, backgroundColor: C.emerald, borderRadius: 4, borderSkipped: false },
  ],
}));

const c4Data = computed(() => ({
  labels: ["Faible","Moyenne","Haute","Critique"],
  datasets: [{ data: [stats.value.low_incidents, stats.value.medium_incidents, stats.value.high_incidents, stats.value.critical_incidents], backgroundColor: [C.slate, C.cyan, C.amber, C.rose], borderRadius: 4, borderSkipped: false }],
}));
const c4Opts = { ...baseOpts, indexAxis: "y", scales: { x: { grid: { color: GRID }, ticks: tickStyle, border: { display: false } }, y: { grid: { display: false }, ticks: tickStyle } } };

const c5Data = computed(() => ({
  labels: ["Critique","Haute","Moyenne","Faible"],
  datasets: [{ data: [stats.value.mttr_critical, stats.value.mttr_high, stats.value.mttr_medium, stats.value.mttr_low], backgroundColor: [C.rose, C.amber, C.cyan, C.slate], borderRadius: 4, borderSkipped: false }],
}));
const c5Opts = { ...baseOpts, scales: { x: { grid: { display: false }, ticks: tickStyle }, y: { grid: { color: GRID }, border: { display: false }, ticks: { ...tickStyle, callback: v => v + "h" } } } };

const c6Data = computed(() => ({
  labels: stats.value.incident_categories?.map(i => i.type) ?? [],
  datasets: [{ data: stats.value.incident_categories?.map(i => i.count) ?? [], backgroundColor: CAT_COLORS, borderWidth: 0, hoverOffset: 5 }],
}));

const c7Data = computed(() => ({
  labels: stats.value.service_workload?.map(i => i.service) ?? [],
  datasets: [
    { label: "Ouverts",  data: stats.value.service_workload?.map(i => i.open)        ?? [], backgroundColor: C.orange, borderRadius: 3, borderSkipped: false },
    { label: "En cours", data: stats.value.service_workload?.map(i => i.in_progress) ?? [], backgroundColor: C.indigo, borderRadius: 3, borderSkipped: false },
  ],
}));

const c8Data = computed(() => ({
  labels: ["S12","S13","S14","S15","S16","S17","S18","S19","S20","S21","S22","S23"],
  datasets: [
    { data: stats.value.sla_history, borderColor: C.indigo, backgroundColor: "rgba(99,102,241,.07)", fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.indigo, borderWidth: 2 },
    { data: Array(12).fill(95), borderColor: C.rose, borderDash: [5,4], borderWidth: 1.5, pointRadius: 0, fill: false },
  ],
}));
const c8Opts = { ...baseOpts, scales: { x: { grid: { display: false }, ticks: { ...tickStyle, autoSkip: false, maxRotation: 0 } }, y: { min: 75, max: 100, grid: { color: GRID }, border: { display: false }, ticks: { ...tickStyle, callback: v => v + "%" } } } };

onMounted(() => { fetchSites(); fetchServices(); fetchStats(); });
</script>

<template>
  <DashboardLayout>
    <div class="db">
      <h2 class="sr-only">Tableau de bord ITSM — métriques d'infrastructure, statuts et graphiques</h2>

      <!-- ══ TOPBAR ══════════════════════════════════════════════════════════ -->
      <header class="topbar">
        <div class="tb-left">
          <div class="tb-env">
            <span class="env-dot" aria-hidden="true" />
            Production
          </div>
          <h1 class="tb-title">Tableau de bord ITSM</h1>
          <p class="tb-sub">
            Infrastructure &amp; incidents
            <span class="dot-sep" aria-hidden="true">·</span>
            Mis à jour à {{ fmtClock(lastUpdated) }}
          </p>
        </div>

        <div class="alert-pills">
          <span v-if="stats.critical_incidents > 0" class="pill p-rose pulse">
            <i class="ti ti-alert-circle" aria-hidden="true" />
            {{ stats.critical_incidents }} critiques
          </span>
          <span v-if="stats.overdue_incidents > 0" class="pill p-amber">
            <i class="ti ti-clock-exclamation" aria-hidden="true" />
            {{ stats.overdue_incidents }} hors SLA
          </span>
          <span v-if="stats.unread_notifications > 0" class="pill p-indigo">
            <i class="ti ti-bell" aria-hidden="true" />
            {{ stats.unread_notifications }}
          </span>
          <span class="pill p-emerald">{{ fmt(stats.sla_resolution_rate) }}% SLA</span>
          <span v-if="isLoading" class="spinner" aria-label="Chargement" />
        </div>
      </header>

      <!-- ══ FILTRES ═════════════════════════════════════════════════════════ -->
      <div class="fbar" role="toolbar" aria-label="Filtres du tableau de bord">

        <div class="fg">
          <label class="flbl" for="f-site">
            <i class="ti ti-map-pin" aria-hidden="true" />Site
          </label>
          <div class="sel-wrap">
            <select id="f-site" v-model="filters.site_id" class="fsel">
              <option value="">Tous les sites</option>
              <option v-for="s in sites" :key="s.ID" :value="s.ID">{{ s.name }}</option>
            </select>
            <i class="ti ti-chevron-down sel-arr" aria-hidden="true" />
          </div>
        </div>

        <span class="fsep" aria-hidden="true" />

        <div class="fg">
          <label class="flbl" for="f-service">
            <i class="ti ti-server" aria-hidden="true" />Service
          </label>
          <div class="sel-wrap">
            <select id="f-service" v-model="filters.service_id" class="fsel">
              <option value="">Tous les services</option>
              <option v-for="s in services" :key="s.ID" :value="s.ID">{{ s.name }}</option>
            </select>
            <i class="ti ti-chevron-down sel-arr" aria-hidden="true" />
          </div>
        </div>

        <span class="fsep" aria-hidden="true" />

        <div class="fg">
          <span class="flbl"><i class="ti ti-calendar" aria-hidden="true" />Période</span>
          <div class="fpills">
            <button
                v-for="o in PERIOD_OPTS" :key="o.value"
                class="fp" :class="{ active: filters.period === o.value }"
                :aria-pressed="filters.period === o.value"
                @click="filters.period = o.value"
            >{{ o.label }}</button>
          </div>
        </div>

      </div>

      <!-- ══ INFRASTRUCTURE ══════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-infra">
        <p id="sec-infra" class="sec">Infrastructure</p>
        <div class="krow">

          <div class="kc">
            <div class="kc-ico" style="background:rgba(99,102,241,.1);color:#6366f1;">
              <i class="ti ti-server-2" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Parc total</div>
            <div class="kc-val">{{ stats.total_assets }}</div>
            <div class="kc-sub">{{ stats.active_assets }} actifs · {{ stats.alert_assets }} alertes · {{ stats.offline_assets }} hors ligne</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: activeAssetsPercent + '%', background: '#6366f1' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(16,185,129,.1);color:#10b981;">
              <i class="ti ti-circle-check" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Actifs</div>
            <div class="kc-val">{{ stats.active_assets }}</div>
            <div class="kc-sub">{{ activeAssetsPercent }}% disponibles</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: activeAssetsPercent + '%', background: '#10b981' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(245,158,11,.1);color:#f59e0b;">
              <i class="ti ti-alert-triangle" aria-hidden="true" />
            </div>
            <div class="kc-lbl">En alerte</div>
            <div class="kc-val">{{ stats.alert_assets }}</div>
            <div class="kc-sub">{{ alertAssetsPercent }}% nécessitent attention</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: alertAssetsPercent + '%', background: '#f59e0b' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(239,68,68,.1);color:#ef4444;">
              <i class="ti ti-plug-off" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Hors ligne</div>
            <div class="kc-val">{{ stats.offline_assets }}</div>
            <div class="kc-sub">{{ offlineAssetsPercent }}% indisponibles</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: offlineAssetsPercent + '%', background: '#ef4444' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(148,163,184,.15);color:#64748b;">
              <i class="ti ti-ticket" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Incidents (filtrés)</div>
            <div class="kc-val">{{ stats.total_incidents }}</div>
            <div class="kc-sub">{{ stats.open_incidents }} ouverts · {{ stats.in_progress_incidents }} en cours</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: resolvedPct + '%', background: '#6366f1' }" /></div>
          </div>

        </div>
      </section>

      <!-- ══ STATUTS ═════════════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-status">
        <p id="sec-status" class="sec">Statuts des incidents</p>
        <div class="srow">

          <div class="sc sc-or">
            <div class="sc-hd">
              <span class="sc-lbl">Ouverts</span>
              <span class="sc-bdg bdg-or">À traiter</span>
            </div>
            <div class="sc-val">{{ stats.open_incidents }}</div>
            <div class="sc-tr"><i class="ti ti-arrow-up-right" aria-hidden="true" />Nouveaux tickets</div>
          </div>

          <div class="sc sc-in">
            <div class="sc-hd">
              <span class="sc-lbl">En cours</span>
              <span class="sc-bdg bdg-in">Assignés</span>
            </div>
            <div class="sc-val">{{ stats.in_progress_incidents }}</div>
            <div class="sc-tr neutral"><i class="ti ti-user-check" aria-hidden="true" />En traitement</div>
          </div>

          <div class="sc sc-em">
            <div class="sc-hd">
              <span class="sc-lbl">Résolus</span>
              <span class="sc-bdg bdg-em">À valider</span>
            </div>
            <div class="sc-val">{{ stats.resolved_incidents }}</div>
            <div class="sc-tr positive"><i class="ti ti-circle-check" aria-hidden="true" />En attente clôture</div>
          </div>

          <div class="sc sc-sl">
            <div class="sc-hd">
              <span class="sc-lbl">Clôturés</span>
              <span class="sc-bdg bdg-sl">Archivés</span>
            </div>
            <div class="sc-val">{{ stats.closed_incidents }}</div>
            <div class="sc-tr neutral"><i class="ti ti-archive" aria-hidden="true" />Total historique</div>
          </div>

          <div class="sc sc-rs" :class="{ 'sc-alert': stats.overdue_incidents > 0 }">
            <div class="sc-hd">
              <span class="sc-lbl">Hors SLA</span>
              <span v-if="stats.overdue_incidents > 0" class="sc-bdg bdg-rs pulse">En retard</span>
              <span v-else class="sc-bdg bdg-em">À jour</span>
            </div>
            <div class="sc-val" :class="{ 'val-rose': stats.overdue_incidents > 0 }">{{ stats.overdue_incidents }}</div>
            <div class="sc-tr" :class="stats.overdue_incidents > 0 ? 'danger' : 'positive'">
              <i class="ti ti-clock-exclamation" aria-hidden="true" />
              {{ stats.overdue_incidents > 0 ? 'Action requise' : 'Objectif atteint' }}
            </div>
          </div>

        </div>
      </section>

      <!-- ══ CMDB ════════════════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-cmdb">
        <p id="sec-cmdb" class="sec">CMDB &amp; Dépendances</p>
        <div class="krow krow-3">

          <div class="kc">
            <div class="kc-ico" style="background:rgba(99,102,241,.1);color:#6366f1;">
              <i class="ti ti-apps" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Applications</div>
            <div class="kc-val">{{ stats.total_applications }}</div>
            <div class="kc-sub">{{ stats.critical_applications }} critiques</div>
            <div class="kc-bar"><div class="kc-fill" style="width:100%;background:#6366f1;" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(6,182,212,.1);color:#06b6d4;">
              <i class="ti ti-database" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Bases de données</div>
            <div class="kc-val">{{ stats.total_databases }}</div>
            <div class="kc-sub">{{ stats.production_databases }} en production</div>
            <div class="kc-bar"><div class="kc-fill" style="width:100%;background:#06b6d4;" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico" style="background:rgba(139,92,246,.1);color:#8b5cf6;">
              <i class="ti ti-hierarchy" aria-hidden="true" />
            </div>
            <div class="kc-lbl">Relations CMDB</div>
            <div class="kc-val">{{ stats.total_relations }}</div>
            <div class="kc-sub">Dépendances cartographiées</div>
            <div class="kc-bar"><div class="kc-fill" style="width:100%;background:#8b5cf6;" /></div>
          </div>

        </div>
      </section>

      <!-- ══ GRAPHIQUES ══════════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-charts">
        <p id="sec-charts" class="sec">Analyse graphique</p>

        <div class="g2">
          <div class="card">
            <div class="card-hd">
              <div><h3>Tendance — incidents</h3><p>Ouverts vs résolus sur la période</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#6366f1;" />Ouverts</span>
                <span class="li"><span class="ld" style="background:#10b981;" />Résolus</span>
              </div>
            </div>
            <div class="cw"><Line :data="c1Data" :options="{ ...baseOpts, scales: commonScales }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Répartition par statut</h3><p>Distribution des {{ stats.total_incidents }} incidents</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#f97316;" />Ouverts {{ stats.open_incidents }}</span>
                <span class="li"><span class="ld" style="background:#6366f1;" />En cours {{ stats.in_progress_incidents }}</span>
                <span class="li"><span class="ld" style="background:#10b981;" />Résolus {{ stats.resolved_incidents }}</span>
                <span class="li"><span class="ld" style="background:#94a3b8;" />Clôturés {{ stats.closed_incidents }}</span>
              </div>
            </div>
            <div class="cw"><Doughnut :data="c2Data" :options="{ ...baseOpts, cutout: '68%' }" /></div>
          </div>
        </div>

        <div class="g2">
          <div class="card">
            <div class="card-hd">
              <div><h3>Volume mensuel</h3><p>Incidents créés vs résolus par mois</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#6366f1;" />Créés</span>
                <span class="li"><span class="ld" style="background:#10b981;" />Résolus</span>
              </div>
            </div>
            <div class="cw"><Bar :data="c3Data" :options="{ ...baseOpts, scales: commonScales }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Priorité des incidents actifs</h3><p>Répartition non clôturée</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#94a3b8;" />Faible {{ stats.low_incidents }}</span>
                <span class="li"><span class="ld" style="background:#06b6d4;" />Moyenne {{ stats.medium_incidents }}</span>
                <span class="li"><span class="ld" style="background:#f59e0b;" />Haute {{ stats.high_incidents }}</span>
                <span class="li"><span class="ld" style="background:#ef4444;" />Critique {{ stats.critical_incidents }}</span>
              </div>
            </div>
            <div class="cw"><Bar :data="c4Data" :options="c4Opts" /></div>
          </div>
        </div>

        <div class="g3">
          <div class="card">
            <div class="card-hd"><div><h3>MTTR par priorité</h3><p>Temps moyen (heures)</p></div></div>
            <div class="cw cw-sm"><Bar :data="c5Data" :options="c5Opts" /></div>
          </div>
          <div class="card">
            <div class="card-hd"><div><h3>Incidents par catégorie</h3><p>Nature des tickets ce mois</p></div></div>
            <div class="leg leg-block">
              <span v-for="(cat, i) in stats.incident_categories" :key="cat.type" class="li">
                <span class="ld" :style="{ background: CAT_COLORS[i % CAT_COLORS.length] }" />
                {{ cat.type }} {{ cat.count }}
              </span>
            </div>
            <div class="cw cw-sm"><Doughnut :data="c6Data" :options="{ ...baseOpts, cutout: '60%' }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Charge par service</h3><p>Tickets ouverts et en cours</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#f97316;" />Ouverts</span>
                <span class="li"><span class="ld" style="background:#6366f1;" />En cours</span>
              </div>
            </div>
            <div class="cw cw-sm"><Bar :data="c7Data" :options="{ ...baseOpts, scales: stackedScales }" /></div>
          </div>
        </div>

        <div class="g1">
          <div class="card">
            <div class="card-hd">
              <div><h3>Évolution SLA — 12 semaines</h3><p>Taux de respect SLA (%) · ligne = objectif 95%</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#6366f1;" />Taux réel</span>
                <span class="li"><span class="ld" style="background:#ef4444;width:16px;height:2px;border-radius:0;" />Objectif 95%</span>
              </div>
            </div>
            <div class="cw"><Line :data="c8Data" :options="c8Opts" /></div>
          </div>
        </div>
      </section>

      <!-- ══ PERFORMANCE SLA ═════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-sla">
        <p id="sec-sla" class="sec">Performance SLA</p>
        <div class="sla-grid">

          <div class="sla-c">
            <div class="sla-hd">
              <div class="sla-ico" :style="{ background: slaColor + '18', color: slaColor }">
                <i class="ti ti-target" aria-hidden="true" />
              </div>
              <div>
                <div class="sla-lbl">Taux résolution SLA</div>
                <div class="sla-val" :style="{ color: slaColor }">{{ fmt(stats.sla_resolution_rate) }}%</div>
              </div>
            </div>
            <div class="sla-trk"><div class="sla-fill" :style="{ width: stats.sla_resolution_rate + '%', background: slaColor }" /></div>
            <div class="sla-ft">
              <span>Objectif 95%</span>
              <span :style="{ color: slaColor }">{{ stats.sla_resolution_rate >= 95 ? '✓ Atteint' : `${fmt(95 - stats.sla_resolution_rate)} pts manquants` }}</span>
            </div>
          </div>

          <div class="sla-c">
            <div class="sla-hd">
              <div class="sla-ico" :style="{ background: mttrColor + '18', color: mttrColor }">
                <i class="ti ti-clock" aria-hidden="true" />
              </div>
              <div>
                <div class="sla-lbl">MTTR moyen</div>
                <div class="sla-val" :style="{ color: mttrColor }">{{ fmtTime(stats.mttr_hours) }}</div>
              </div>
            </div>
            <div class="sla-trk"><div class="sla-fill" :style="{ width: Math.min((stats.mttr_hours / 6) * 100, 100) + '%', background: mttrColor }" /></div>
            <div class="sla-ft">
              <span>Cible maximale 3h</span>
              <span :style="{ color: mttrColor }">{{ stats.mttr_hours <= 3 ? "✓ Dans l'objectif" : 'Hors objectif' }}</span>
            </div>
          </div>

          <div class="sla-c">
            <div class="sla-hd">
              <div class="sla-ico" style="background:rgba(16,185,129,.1);color:#10b981;">
                <i class="ti ti-refresh" aria-hidden="true" />
              </div>
              <div>
                <div class="sla-lbl">Taux réouverture</div>
                <div class="sla-val" :style="{ color: stats.reopening_rate < 5 ? '#10b981' : '#ef4444' }">{{ fmt(stats.reopening_rate) }}%</div>
              </div>
            </div>
            <div class="sla-trk"><div class="sla-fill" :style="{ width: (100 - stats.reopening_rate) + '%', background: stats.reopening_rate < 5 ? '#10b981' : '#ef4444' }" /></div>
            <div class="sla-ft">
              <span>Objectif &lt; 5%</span>
              <span :style="{ color: stats.reopening_rate < 5 ? '#10b981' : '#ef4444' }">{{ stats.reopening_rate < 5 ? "✓ Dans l'objectif" : 'À surveiller' }}</span>
            </div>
          </div>

        </div>
      </section>
    </div>
  </DashboardLayout>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&display=swap');

.db { padding: 0 0 2rem; font-family: 'DM Sans', var(--font-sans); }

/* ── Topbar ───────────────────────────────────────────────────────────────── */
.topbar { display: flex; justify-content: space-between; align-items: flex-end; padding: 1.5rem 0 1.25rem; border-bottom: 0.5px solid var(--color-border-tertiary); margin-bottom: 0; }
.tb-env { display: flex; align-items: center; gap: 6px; font-size: 11px; font-weight: 500; letter-spacing: .08em; text-transform: uppercase; color: #10b981; margin-bottom: 5px; }
.env-dot { width: 7px; height: 7px; border-radius: 50%; background: #10b981; }
.tb-title { font-size: 20px; font-weight: 500; color: var(--color-text-primary); letter-spacing: -.4px; line-height: 1.15; }
.tb-sub { font-size: 12px; color: var(--color-text-tertiary); margin-top: 3px; }
.dot-sep { margin: 0 6px; opacity: .4; }

/* Alert pills */
.alert-pills { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; justify-content: flex-end; }
.pill { display: inline-flex; align-items: center; gap: 4px; font-size: 11px; padding: 4px 10px; border-radius: 100px; border: 0.5px solid transparent; }
.p-rose    { background: #FCEBEB; color: #A32D2D; border-color: #F09595; }
.p-amber   { background: #FAEEDA; color: #854F0B; border-color: #FAC775; }
.p-indigo  { background: #EEEDFE; color: #3C3489; border-color: #AFA9EC; }
.p-emerald { background: #EAF3DE; color: #3B6D11; border-color: #C0DD97; }

/* ── Filter bar ───────────────────────────────────────────────────────────── */
.fbar { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; padding: 10px 0; border-bottom: 0.5px solid var(--color-border-tertiary); margin-bottom: 1.5rem; }
.fg   { display: flex; align-items: center; gap: 8px; }
.flbl { font-size: 10px; font-weight: 500; letter-spacing: .08em; text-transform: uppercase; color: var(--color-text-tertiary); white-space: nowrap; display: flex; align-items: center; gap: 4px; }
.flbl .ti { font-size: 12px; }
.fsep { width: 0.5px; height: 20px; background: var(--color-border-tertiary); flex-shrink: 0; }

.sel-wrap { position: relative; }
.fsel { padding: 5px 26px 5px 10px; font-size: 12px; font-family: 'DM Sans', var(--font-sans); background: var(--color-background-primary); border: 0.5px solid var(--color-border-tertiary); border-radius: var(--border-radius-md); color: var(--color-text-primary); appearance: none; cursor: pointer; outline: none; transition: border-color .15s; }
.fsel:focus { border-color: #6366f1; }
.sel-arr { position: absolute; right: 8px; top: 50%; transform: translateY(-50%); font-size: 11px; color: var(--color-text-tertiary); pointer-events: none; }

.fpills { display: flex; gap: 4px; }
.fp { padding: 4px 10px; font-size: 11px; font-family: 'DM Sans', var(--font-sans); border-radius: 100px; border: 0.5px solid var(--color-border-tertiary); background: transparent; color: var(--color-text-secondary); cursor: pointer; transition: background .12s, border-color .12s, color .12s; }
.fp:hover  { background: var(--color-background-secondary); }
.fp.active { background: #6366f1; color: #fff; border-color: #6366f1; }

/* ── Section title ────────────────────────────────────────────────────────── */
.sec { font-size: 10px; font-weight: 500; letter-spacing: .1em; text-transform: uppercase; color: var(--color-text-tertiary); margin: 1.5rem 0 .75rem; }

/* ── KPI cards ────────────────────────────────────────────────────────────── */
.krow   { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; }
.krow-3 { grid-template-columns: repeat(3, 1fr); }
.kc { background: var(--color-background-primary); border: 0.5px solid var(--color-border-tertiary); border-radius: var(--border-radius-lg); padding: 14px; transition: border-color .15s; }
.kc:hover { border-color: var(--color-border-secondary); }
.kc-ico { width: 32px; height: 32px; border-radius: var(--border-radius-md); display: flex; align-items: center; justify-content: center; font-size: 15px; margin-bottom: 10px; }
.kc-lbl { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 4px; }
.kc-val { font-size: 26px; font-weight: 500; color: var(--color-text-primary); letter-spacing: -1.2px; line-height: 1; font-variant-numeric: tabular-nums; }
.kc-sub { font-size: 11px; color: var(--color-text-tertiary); margin-top: 5px; }
.kc-bar { height: 2px; background: var(--color-border-tertiary); border-radius: 2px; margin-top: 12px; overflow: hidden; }
.kc-fill { height: 100%; border-radius: 2px; transition: width .4s ease; }

/* ── Status cards ─────────────────────────────────────────────────────────── */
.srow { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; }
.sc { background: var(--color-background-primary); border: 0.5px solid var(--color-border-tertiary); border-radius: var(--border-radius-lg); padding: 14px; position: relative; overflow: hidden; transition: border-color .15s; }
.sc::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; border-radius: 0; }
.sc-or::before { background: #f97316; }
.sc-in::before { background: #6366f1; }
.sc-em::before { background: #10b981; }
.sc-sl::before { background: #94a3b8; }
.sc-rs::before { background: #ef4444; }
.sc-alert { background: rgba(239,68,68,.03); border-color: rgba(239,68,68,.2); }

.sc-hd { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.sc-lbl { font-size: 11px; font-weight: 500; color: var(--color-text-secondary); }
.sc-bdg { font-size: 9px; font-weight: 500; padding: 2px 7px; border-radius: 100px; text-transform: uppercase; letter-spacing: .04em; }
.bdg-or { background: #FAECE7; color: #993C1D; }
.bdg-in { background: #EEEDFE; color: #3C3489; }
.bdg-em { background: #EAF3DE; color: #3B6D11; }
.bdg-sl { background: #F1EFE8; color: #5F5E5A; }
.bdg-rs { background: #FCEBEB; color: #A32D2D; }

.sc-val { font-size: 30px; font-weight: 500; color: var(--color-text-primary); letter-spacing: -1.8px; line-height: 1; font-variant-numeric: tabular-nums; }
.val-rose { color: #ef4444; }
.sc-tr { display: flex; align-items: center; gap: 4px; margin-top: 8px; font-size: 11px; color: var(--color-text-tertiary); }
.sc-tr.positive { color: #10b981; }
.sc-tr.danger   { color: #ef4444; }
.sc-tr.neutral  { color: var(--color-text-tertiary); }

/* ── Chart grids ──────────────────────────────────────────────────────────── */
.g1 { display: grid; grid-template-columns: 1fr;           gap: 10px; margin-bottom: 10px; }
.g2 { display: grid; grid-template-columns: 1fr 1fr;       gap: 10px; margin-bottom: 10px; }
.g3 { display: grid; grid-template-columns: 1fr 1fr 1fr;   gap: 10px; margin-bottom: 10px; }

/* ── Chart card ───────────────────────────────────────────────────────────── */
.card { background: var(--color-background-primary); border: 0.5px solid var(--color-border-tertiary); border-radius: var(--border-radius-lg); padding: 16px; transition: border-color .15s; }
.card:hover { border-color: var(--color-border-secondary); }
.card-hd { display: flex; justify-content: space-between; align-items: flex-start; gap: 12px; margin-bottom: 12px; }
.card h3 { font-size: 13px; font-weight: 500; color: var(--color-text-primary); margin-bottom: 2px; }
.card p  { font-size: 11px; color: var(--color-text-secondary); }
.cw     { position: relative; height: 175px; }
.cw-sm  { height: 155px; }

/* ── Legend ───────────────────────────────────────────────────────────────── */
.leg       { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; flex-shrink: 0; }
.leg-block { margin-bottom: 8px; }
.li  { display: flex; align-items: center; gap: 5px; font-size: 11px; color: var(--color-text-secondary); }
.ld  { width: 8px; height: 8px; border-radius: 2px; flex-shrink: 0; }

/* ── SLA cards ────────────────────────────────────────────────────────────── */
.sla-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-bottom: 1.25rem; }
.sla-c { background: var(--color-background-primary); border: 0.5px solid var(--color-border-tertiary); border-radius: var(--border-radius-lg); padding: 16px; transition: border-color .15s; }
.sla-c:hover { border-color: var(--color-border-secondary); }
.sla-hd { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.sla-ico { width: 36px; height: 36px; border-radius: var(--border-radius-md); display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; }
.sla-lbl { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 3px; }
.sla-val { font-size: 22px; font-weight: 500; letter-spacing: -.5px; font-variant-numeric: tabular-nums; }
.sla-trk { height: 4px; background: var(--color-border-tertiary); border-radius: 2px; overflow: hidden; }
.sla-fill { height: 100%; border-radius: 2px; transition: width .5s ease; }
.sla-ft { display: flex; justify-content: space-between; font-size: 11px; color: var(--color-text-tertiary); margin-top: 7px; }

/* ── Spinner & animations ─────────────────────────────────────────────────── */
.spinner { width: 14px; height: 14px; border: 1.5px solid var(--color-border-tertiary); border-top-color: #6366f1; border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin  { to { transform: rotate(360deg); } }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: .4; } }
.pulse { animation: pulse 1.8s infinite; }
</style>


<script setup>
import { computed, onMounted, ref, watch } from "vue";
import {
  Chart as ChartJS,
  ArcElement, BarElement, CategoryScale, LinearScale,
  PointElement, LineElement, Tooltip, Legend, Filler,
} from "chart.js";
import { Doughnut, Bar, Line } from "vue-chartjs";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

ChartJS.register(
    ArcElement, BarElement, CategoryScale, LinearScale,
    PointElement, LineElement, Tooltip, Legend, Filler
);

// ─── Palette — accordée au thème dark du projet ───────────────────────────────
// Les couleurs des graphiques sont hardcodées (Chart.js ne lit pas les CSS vars)
// mais restent lisibles sur fond sombre : teintes mid-range (#7F77DD, #1D9E75…)
const C = {
  brand:   "#7f77dd",   // violet — couleur primaire de l'app
  teal:    "#1d9e75",   // vert — résolu / succès
  coral:   "#d85a30",   // orange — ouvert / alerte
  blue:    "#378add",   // bleu — info
  amber:   "#ef9f27",   // amber — warning / priorité haute
  red:     "#e24b4a",   // rouge — danger / critique
  gray:    "#888780",   // gris — neutre / clôturé
};

// Couleurs des catégories d'incidents (ordre fixe)
const CATEGORY_COLORS = [C.brand, C.blue, C.teal, C.amber, C.gray];

// Config Chart.js — axes et grille adaptés au fond sombre
const CHART = {
  tick: { size: 10, color: "#64748b" },
  grid: "rgba(71,85,105,0.25)",        // plus visible sur fond sombre que 0.08
};

// ─── Filtres ──────────────────────────────────────────────────────────────────
const sites    = ref([]);
const services = ref([]);

const filters  = ref({ site_id: "", service_id: "", period: "30" });

const PERIOD_OPTS = [
  { label: "7 j",    value: "7"   },
  { label: "30 j",   value: "30"  },
  { label: "90 j",   value: "90"  },
  { label: "6 mois", value: "180" },
];

// ─── État ─────────────────────────────────────────────────────────────────────
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
    stats.value = data;
    lastUpdated.value = new Date();
  } catch (err) {
    console.error("[Dashboard] fetchStats:", err);
  } finally {
    isLoading.value = false;
  }
};

// Recharge automatiquement à chaque changement de filtre
watch(filters, fetchStats, { deep: true });

// ─── Computed ─────────────────────────────────────────────────────────────────
const pct = (k) => computed(() =>
    stats.value.total_assets
        ? Math.round((stats.value[k] / stats.value.total_assets) * 100) : 0
);
const activeAssetsPercent  = pct("active_assets");
const alertAssetsPercent   = pct("alert_assets");
const offlineAssetsPercent = pct("offline_assets");

const resolvedIncidentPercent = computed(() =>
    stats.value.total_incidents
        ? Math.round(((stats.value.resolved_incidents + stats.value.closed_incidents)
            / stats.value.total_incidents) * 100) : 0
);

// Couleur dynamique de la jauge SLA
const slaColor = computed(() =>
    stats.value.sla_resolution_rate >= 90 ? C.teal
        : stats.value.sla_resolution_rate >= 75 ? C.amber : C.red
);
// Couleur dynamique du MTTR
const mttrColor = computed(() =>
    stats.value.mttr_hours <= 3 ? C.teal
        : stats.value.mttr_hours <= 6 ? C.amber : C.red
);

// ─── Helpers ──────────────────────────────────────────────────────────────────
const fmt     = (v, d = 1) => v != null ? Number(v).toFixed(d) : "0";
const fmtTime = (h) => {
  if (!h) return "—";
  const hrs = Math.floor(h);
  const m   = Math.round((h - hrs) * 60);
  return m ? `${hrs}h ${m}m` : `${hrs}h`;
};
const fmtClock = (d) => d.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" });

// ─── Chart.js options partagées ───────────────────────────────────────────────
const baseOpts = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  layout: { padding: { top: 4, bottom: 4 } },
};

const ticks = CHART.tick;

const commonScales = {
  x: { grid: { display: false }, ticks },
  y: { grid: { color: CHART.grid }, ticks, border: { display: false } },
};
const stackedScales = {
  x: { stacked: true, grid: { display: false }, ticks },
  y: { stacked: true, grid: { color: CHART.grid }, ticks, border: { display: false } },
};

// ─── Datasets ─────────────────────────────────────────────────────────────────

/** C1 — Courbe tendance ouverts vs résolus */
const c1Data = computed(() => ({
  labels: ["Lun","Mar","Mer","Jeu","Ven","Sam","Dim"],
  datasets: [
    { data: stats.value.weekly_trend_opened,   borderColor: C.brand, backgroundColor: "rgba(127,119,221,.1)",  fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.brand, borderWidth: 2 },
    { data: stats.value.weekly_trend_resolved, borderColor: C.teal,  backgroundColor: "rgba(29,158,117,.08)", fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.teal,  borderWidth: 2 },
  ],
}));

/** C2 — Donut répartition par statut */
const c2Data = computed(() => ({
  labels: ["Ouverts","En cours","Résolus","Clôturés"],
  datasets: [{
    data: [stats.value.open_incidents, stats.value.in_progress_incidents,
      stats.value.resolved_incidents, stats.value.closed_incidents],
    backgroundColor: [C.coral, C.brand, C.teal, C.gray],
    borderWidth: 0, hoverOffset: 6,
  }],
}));

/** C3 — Barres groupées volume mensuel */
const c3Data = computed(() => ({
  labels: ["Jan","Fév","Mar","Avr","Mai","Jun"],
  datasets: [
    { data: stats.value.monthly_created,  backgroundColor: C.brand, borderRadius: 4, borderSkipped: false },
    { data: stats.value.monthly_resolved, backgroundColor: C.teal,  borderRadius: 4, borderSkipped: false },
  ],
}));

/** C4 — Barres horizontales priorités */
const c4Data = computed(() => ({
  labels: ["Faible","Moyenne","Haute","Critique"],
  datasets: [{
    data: [stats.value.low_incidents, stats.value.medium_incidents,
      stats.value.high_incidents, stats.value.critical_incidents],
    backgroundColor: [C.gray, C.blue, C.amber, C.red],
    borderRadius: 4, borderSkipped: false,
  }],
}));
const c4Opts = {
  ...baseOpts, indexAxis: "y",
  scales: {
    x: { grid: { color: CHART.grid }, ticks, border: { display: false } },
    y: { grid: { display: false }, ticks },
  },
};

/** C5 — MTTR par priorité */
const c5Data = computed(() => ({
  labels: ["Critique","Haute","Moyenne","Faible"],
  datasets: [{
    data: [stats.value.mttr_critical, stats.value.mttr_high,
      stats.value.mttr_medium,   stats.value.mttr_low],
    backgroundColor: [C.red, C.amber, C.blue, C.gray],
    borderRadius: 4, borderSkipped: false,
  }],
}));
const c5Opts = {
  ...baseOpts,
  scales: {
    x: { grid: { display: false }, ticks },
    y: { grid: { color: CHART.grid }, border: { display: false },
      ticks: { ...ticks, callback: v => v + "h" } },
  },
};

/** C6 — Donut catégories */
const c6Data = computed(() => ({
  labels: stats.value.incident_categories?.map(i => i.type) ?? [],
  datasets: [{
    data: stats.value.incident_categories?.map(i => i.count) ?? [],
    backgroundColor: CATEGORY_COLORS,
    borderWidth: 0, hoverOffset: 5,
  }],
}));

/** C7 — Barres empilées charge service */
const c7Data = computed(() => ({
  labels: stats.value.service_workload?.map(i => i.service) ?? [],
  datasets: [
    { label: "Ouverts",  data: stats.value.service_workload?.map(i => i.open)        ?? [], backgroundColor: C.coral, borderRadius: 3, borderSkipped: false },
    { label: "En cours", data: stats.value.service_workload?.map(i => i.in_progress) ?? [], backgroundColor: C.brand, borderRadius: 3, borderSkipped: false },
  ],
}));

/** C8 — Courbe SLA + objectif 95% */
const c8Data = computed(() => ({
  labels: ["S12","S13","S14","S15","S16","S17","S18","S19","S20","S21","S22","S23"],
  datasets: [
    { data: stats.value.sla_history, borderColor: C.brand, backgroundColor: "rgba(127,119,221,.1)", fill: true, tension: .4, pointRadius: 3, pointBackgroundColor: C.brand, borderWidth: 2 },
    { data: Array(12).fill(95), borderColor: C.red, borderDash: [5,4], borderWidth: 1.5, pointRadius: 0, fill: false },
  ],
}));
const c8Opts = {
  ...baseOpts,
  scales: {
    x: { grid: { display: false }, ticks: { ...ticks, autoSkip: false, maxRotation: 0 } },
    y: { min: 75, max: 100, grid: { color: CHART.grid }, border: { display: false },
      ticks: { ...ticks, callback: v => v + "%" } },
  },
};

onMounted(() => { fetchSites(); fetchServices(); fetchStats(); });
</script>

<template>
  <DashboardLayout>
    <div class="db">
      <h2 class="sr-only">Tableau de bord ITSM — métriques, tendances et analyse des incidents</h2>

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

        <!-- Filtres intégrés dans la topbar -->
        <div class="tb-filters" role="group" aria-label="Filtres">
          <div class="sel-wrap">
            <select v-model="filters.site_id" class="tb-select" aria-label="Filtrer par site">
              <option value="">Tous les sites</option>
              <option v-for="s in sites" :key="s.ID" :value="s.ID">{{ s.name }}</option>
            </select>
            <i class="ti ti-chevron-down sel-arr" aria-hidden="true" />
          </div>
          <div class="sel-wrap">
            <select v-model="filters.service_id" class="tb-select" aria-label="Filtrer par service">
              <option value="">Tous les services</option>
              <option v-for="s in services" :key="s.ID" :value="s.ID">{{ s.name }}</option>
            </select>
            <i class="ti ti-chevron-down sel-arr" aria-hidden="true" />
          </div>
          <div class="period-pills">
            <button
                v-for="o in PERIOD_OPTS" :key="o.value"
                class="pp" :class="{ active: filters.period === o.value }"
                :aria-pressed="filters.period === o.value"
                @click="filters.period = o.value"
            >{{ o.label }}</button>
          </div>
        </div>

        <!-- Badges d'alerte -->
        <div class="alert-pills">
          <span v-if="stats.critical_incidents > 0" class="pill p-red pulse">
            <i class="ti ti-alert-circle" aria-hidden="true" />
            {{ stats.critical_incidents }} critiques
          </span>
          <span v-if="stats.overdue_incidents > 0" class="pill p-amber">
            <i class="ti ti-clock-exclamation" aria-hidden="true" />
            {{ stats.overdue_incidents }} hors SLA
          </span>
          <span v-if="stats.unread_notifications > 0" class="pill p-brand">
            <i class="ti ti-bell" aria-hidden="true" />
            {{ stats.unread_notifications }}
          </span>
          <span class="pill p-green">{{ fmt(stats.sla_resolution_rate) }}% SLA</span>
          <span v-if="isLoading" class="spinner" aria-label="Chargement" />
        </div>
      </header>

      <!-- ══ INFRASTRUCTURE ══════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-infra">
        <p id="sec-infra" class="sec">Infrastructure</p>
        <div class="krow">

          <div class="kc">
            <div class="kc-ico ico-brand"><i class="ti ti-server-2" aria-hidden="true" /></div>
            <div class="kc-lbl">Parc total</div>
            <div class="kc-val">{{ stats.total_assets }}</div>
            <div class="kc-sub">{{ stats.active_assets }} actifs · {{ stats.alert_assets }} alertes · {{ stats.offline_assets }} hors ligne</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: activeAssetsPercent + '%', background: '#7f77dd' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-green"><i class="ti ti-circle-check" aria-hidden="true" /></div>
            <div class="kc-lbl">Actifs</div>
            <div class="kc-val">{{ stats.active_assets }}</div>
            <div class="kc-sub">{{ activeAssetsPercent }}% disponibles</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: activeAssetsPercent + '%', background: '#1d9e75' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-amber"><i class="ti ti-alert-triangle" aria-hidden="true" /></div>
            <div class="kc-lbl">En alerte</div>
            <div class="kc-val">{{ stats.alert_assets }}</div>
            <div class="kc-sub">{{ alertAssetsPercent }}% nécessitent attention</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: alertAssetsPercent + '%', background: '#ef9f27' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-red"><i class="ti ti-plug-off" aria-hidden="true" /></div>
            <div class="kc-lbl">Hors ligne</div>
            <div class="kc-val">{{ stats.offline_assets }}</div>
            <div class="kc-sub">{{ offlineAssetsPercent }}% indisponibles</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: offlineAssetsPercent + '%', background: '#e24b4a' }" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-muted"><i class="ti ti-ticket" aria-hidden="true" /></div>
            <div class="kc-lbl">Incidents filtrés</div>
            <div class="kc-val">{{ stats.total_incidents }}</div>
            <div class="kc-sub">{{ stats.open_incidents }} ouverts · {{ stats.in_progress_incidents }} en cours</div>
            <div class="kc-bar"><div class="kc-fill" :style="{ width: resolvedIncidentPercent + '%', background: '#7f77dd' }" /></div>
          </div>

        </div>
      </section>

      <!-- ══ STATUTS ═════════════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-status">
        <p id="sec-status" class="sec">Statuts des incidents</p>
        <div class="srow">

          <div class="sc sc-coral">
            <div class="sc-hd"><span class="sc-lbl">Ouverts</span><span class="sc-bdg bdg-coral">À traiter</span></div>
            <div class="sc-val">{{ stats.open_incidents }}</div>
            <div class="sc-tr"><i class="ti ti-arrow-up-right" aria-hidden="true" />Nouveaux tickets</div>
          </div>

          <div class="sc sc-brand">
            <div class="sc-hd"><span class="sc-lbl">En cours</span><span class="sc-bdg bdg-brand">Assignés</span></div>
            <div class="sc-val">{{ stats.in_progress_incidents }}</div>
            <div class="sc-tr neutral"><i class="ti ti-user-check" aria-hidden="true" />En traitement</div>
          </div>

          <div class="sc sc-green">
            <div class="sc-hd"><span class="sc-lbl">Résolus</span><span class="sc-bdg bdg-green">À valider</span></div>
            <div class="sc-val">{{ stats.resolved_incidents }}</div>
            <div class="sc-tr positive"><i class="ti ti-circle-check" aria-hidden="true" />En attente clôture</div>
          </div>

          <div class="sc sc-muted">
            <div class="sc-hd"><span class="sc-lbl">Clôturés</span><span class="sc-bdg bdg-muted">Archivés</span></div>
            <div class="sc-val">{{ stats.closed_incidents }}</div>
            <div class="sc-tr neutral"><i class="ti ti-archive" aria-hidden="true" />Total historique</div>
          </div>

          <div class="sc sc-red" :class="{ 'sc-alert': stats.overdue_incidents > 0 }">
            <div class="sc-hd">
              <span class="sc-lbl">Hors SLA</span>
              <span v-if="stats.overdue_incidents > 0" class="sc-bdg bdg-red pulse">En retard</span>
              <span v-else class="sc-bdg bdg-green">À jour</span>
            </div>
            <div class="sc-val" :class="{ 'val-red': stats.overdue_incidents > 0 }">{{ stats.overdue_incidents }}</div>
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
            <div class="kc-ico ico-brand"><i class="ti ti-apps" aria-hidden="true" /></div>
            <div class="kc-lbl">Applications</div>
            <div class="kc-val">{{ stats.total_applications }}</div>
            <div class="kc-sub">{{ stats.critical_applications }} critiques</div>
            <div class="kc-bar"><div class="kc-fill" style="width:100%;background:#7f77dd;" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-blue"><i class="ti ti-database" aria-hidden="true" /></div>
            <div class="kc-lbl">Bases de données</div>
            <div class="kc-val">{{ stats.total_databases }}</div>
            <div class="kc-sub">{{ stats.production_databases }} en production</div>
            <div class="kc-bar"><div class="kc-fill" style="width:100%;background:#378add;" /></div>
          </div>

          <div class="kc">
            <div class="kc-ico ico-violet"><i class="ti ti-hierarchy" aria-hidden="true" /></div>
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
                <span class="li"><span class="ld" style="background:#7f77dd;" />Ouverts</span>
                <span class="li"><span class="ld" style="background:#1d9e75;" />Résolus</span>
              </div>
            </div>
            <div class="cw"><Line :data="c1Data" :options="{ ...baseOpts, scales: commonScales }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Répartition par statut</h3><p>Distribution des {{ stats.total_incidents }} incidents</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#d85a30;" />Ouverts {{ stats.open_incidents }}</span>
                <span class="li"><span class="ld" style="background:#7f77dd;" />En cours {{ stats.in_progress_incidents }}</span>
                <span class="li"><span class="ld" style="background:#1d9e75;" />Résolus {{ stats.resolved_incidents }}</span>
                <span class="li"><span class="ld" style="background:#888780;" />Clôturés {{ stats.closed_incidents }}</span>
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
                <span class="li"><span class="ld" style="background:#7f77dd;" />Créés</span>
                <span class="li"><span class="ld" style="background:#1d9e75;" />Résolus</span>
              </div>
            </div>
            <div class="cw"><Bar :data="c3Data" :options="{ ...baseOpts, scales: commonScales }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Priorité des incidents actifs</h3><p>Répartition non clôturée</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#888780;" />Faible {{ stats.low_incidents }}</span>
                <span class="li"><span class="ld" style="background:#378add;" />Moyenne {{ stats.medium_incidents }}</span>
                <span class="li"><span class="ld" style="background:#ef9f27;" />Haute {{ stats.high_incidents }}</span>
                <span class="li"><span class="ld" style="background:#e24b4a;" />Critique {{ stats.critical_incidents }}</span>
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
            <div class="leg leg-b">
              <span v-for="(cat, i) in stats.incident_categories" :key="cat.type" class="li">
                <span class="ld" :style="{ background: CATEGORY_COLORS[i % CATEGORY_COLORS.length] }" />
                {{ cat.type }} {{ cat.count }}
              </span>
            </div>
            <div class="cw cw-sm"><Doughnut :data="c6Data" :options="{ ...baseOpts, cutout: '60%' }" /></div>
          </div>
          <div class="card">
            <div class="card-hd">
              <div><h3>Charge par service</h3><p>Tickets ouverts et en cours</p></div>
              <div class="leg">
                <span class="li"><span class="ld" style="background:#d85a30;" />Ouverts</span>
                <span class="li"><span class="ld" style="background:#7f77dd;" />En cours</span>
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
                <span class="li"><span class="ld" style="background:#7f77dd;" />Taux réel</span>
                <span class="li"><span class="ld" style="background:#e24b4a;width:16px;height:2px;border-radius:0;" />Objectif 95%</span>
              </div>
            </div>
            <div class="cw cw-sla"><Line :data="c8Data" :options="c8Opts" /></div>
          </div>
        </div>
      </section>

      <!-- ══ PERFORMANCE SLA ═════════════════════════════════════════════════ -->
      <section aria-labelledby="sec-sla">
        <p id="sec-sla" class="sec">Performance SLA</p>
        <div class="sla-grid">

          <div class="sla-c">
            <div class="sla-hd">
              <div class="sla-ico" :style="{ background: slaColor + '22', color: slaColor }">
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
              <div class="sla-ico" :style="{ background: mttrColor + '22', color: mttrColor }">
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
              <div class="sla-ico" style="background:rgba(29,158,117,.2);color:#1d9e75;">
                <i class="ti ti-refresh" aria-hidden="true" />
              </div>
              <div>
                <div class="sla-lbl">Taux réouverture</div>
                <div class="sla-val" :style="{ color: stats.reopening_rate < 5 ? '#1d9e75' : '#e24b4a' }">{{ fmt(stats.reopening_rate) }}%</div>
              </div>
            </div>
            <div class="sla-trk"><div class="sla-fill" :style="{ width: (100 - stats.reopening_rate) + '%', background: stats.reopening_rate < 5 ? '#1d9e75' : '#e24b4a' }" /></div>
            <div class="sla-ft">
              <span>Objectif &lt; 5%</span>
              <span :style="{ color: stats.reopening_rate < 5 ? '#1d9e75' : '#e24b4a' }">{{ stats.reopening_rate < 5 ? "✓ Dans l'objectif" : 'À surveiller' }}</span>
            </div>
          </div>

        </div>
      </section>
    </div>
  </DashboardLayout>
</template>

<style scoped>


.db { padding: 0 0 2rem; }

/* ── Topbar ───────────────────────────────────────────────────────────────── */
.topbar {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  flex-wrap: wrap;
  padding: 1.25rem 1.5rem;
  border-radius: var(--r-xl, 20px);
  margin-bottom: 1.5rem;
}
.tb-left { flex: 0 0 auto; }
.tb-env {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: .1em;
  text-transform: uppercase;
  color: var(--c-success-t, #34d399);
  margin-bottom: 4px;
}
.env-dot {
  width: 6px; height: 6px;
  border-radius: 50%;
  background: currentColor;
}
.tb-title {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -.3px;
  line-height: 1.2;
}
.tb-sub { font-size: 12px; margin-top: 3px; }
.dot-sep { margin: 0 5px; opacity: .35; }

/* Filtres dans la topbar */
.tb-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  flex: 1;
  justify-content: center;
}
.sel-wrap { position: relative; }
.tb-select {
  padding: 5px 26px 5px 10px;
  font-size: 12px;
  border-radius: var(--r-md, 12px);
  appearance: none;
  cursor: pointer;
  min-width: 140px;
}
.sel-arr {
  position: absolute;
  right: 8px; top: 50%;
  transform: translateY(-50%);
  font-size: 11px;
  pointer-events: none;
  opacity: .5;
}
.period-pills { display: flex; gap: 3px; }
.pp {
  padding: 4px 10px;
  font-size: 11px;
  font-family: inherit;
  border-radius: 100px;
  cursor: pointer;
  transition: background .12s, color .12s;
}
.pp.active { background: var(--brand, #7f77dd); color: #fff; border-color: var(--brand, #7f77dd); }

/* Alert pills */
.alert-pills { display: flex; align-items: center; gap: 6px; flex-wrap: wrap; margin-left: auto; }
.pill { display: inline-flex; align-items: center; gap: 4px; font-size: 11px; padding: 4px 10px; border-radius: 100px; }
.p-red   { background: rgba(226,75,74,.15);   color: #f87171; border: 1px solid rgba(226,75,74,.3); }
.p-amber { background: rgba(239,159,39,.12);  color: #fbbf24; border: 1px solid rgba(239,159,39,.25); }
.p-brand { background: rgba(127,119,221,.15); color: #a5b4fc; border: 1px solid rgba(127,119,221,.3); }
.p-green { background: rgba(29,158,117,.15);  color: #34d399; border: 1px solid rgba(29,158,117,.3); }

/* ── Section title ────────────────────────────────────────────────────────── */
.sec {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: .1em;
  text-transform: uppercase;
  margin: 1.5rem 0 .75rem;
}

/* ── KPI cards ────────────────────────────────────────────────────────────── */
.krow   { display: grid; grid-template-columns: repeat(5, 1fr); gap: 10px; }
.krow-3 { grid-template-columns: repeat(3, 1fr); }
.kc {
  border-radius: var(--r-lg, 16px);
  padding: 14px;
  transition: border-color .15s;
}
.kc-ico {
  width: 34px; height: 34px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  margin-bottom: 10px;
}
/* Variantes d'icônes — fonds semi-transparents sur fond sombre */
.ico-brand  { background: rgba(127,119,221,.18); color: #a5b4fc; }
.ico-green  { background: rgba(29,158,117,.18);  color: #34d399; }
.ico-amber  { background: rgba(239,159,39,.18);  color: #fbbf24; }
.ico-red    { background: rgba(226,75,74,.18);   color: #f87171; }
.ico-blue   { background: rgba(55,138,221,.18);  color: #60a5fa; }
.ico-violet { background: rgba(139,92,246,.18);  color: #c084fc; }
.ico-muted  { background: rgba(148,163,184,.12); color: #94a3b8; }

.kc-lbl { font-size: 11px; margin-bottom: 4px; }
.kc-val { font-size: 26px; font-weight: 600; letter-spacing: -1.2px; line-height: 1; font-variant-numeric: tabular-nums; }
.kc-sub { font-size: 11px; margin-top: 5px; }
.kc-bar { height: 2px; border-radius: 2px; margin-top: 12px; overflow: hidden; }
.kc-fill { height: 100%; border-radius: 2px; transition: width .4s ease; }

/* ── Status cards ─────────────────────────────────────────────────────────── */
.srow { display: grid; grid-template-columns: repeat(5, 1fr); gap: 10px; }
.sc {
  border-radius: var(--r-lg, 16px);
  padding: 14px;
  position: relative;
  overflow: hidden;
  transition: border-color .15s;
}
/* Barre colorée en haut = identifiant visuel du statut */
.sc::before { content: ''; position: absolute; top: 0; left: 0; right: 0; height: 2px; border-radius: 0; }
.sc-coral::before { background: #d85a30; }
.sc-brand::before { background: #7f77dd; }
.sc-green::before { background: #1d9e75; }
.sc-muted::before { background: #888780; }
.sc-red::before   { background: #e24b4a; }
.sc-alert         { border-color: rgba(226,75,74,.3) !important; }

.sc-hd  { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.sc-lbl { font-size: 11px; font-weight: 500; }
.sc-bdg { font-size: 9px; font-weight: 600; padding: 2px 7px; border-radius: 100px; text-transform: uppercase; letter-spacing: .04em; }
.bdg-coral { background: rgba(216,90,48,.18);  color: #f97316; }
.bdg-brand { background: rgba(127,119,221,.18); color: #a5b4fc; }
.bdg-green { background: rgba(29,158,117,.18);  color: #34d399; }
.bdg-muted { background: rgba(148,163,184,.12); color: #94a3b8; }
.bdg-red   { background: rgba(226,75,74,.18);   color: #f87171; }

.sc-val { font-size: 28px; font-weight: 600; letter-spacing: -1.5px; line-height: 1; font-variant-numeric: tabular-nums; }
.val-red { color: #f87171; }
.sc-tr  { display: flex; align-items: center; gap: 4px; margin-top: 8px; font-size: 11px; }
.sc-tr.positive { color: #34d399; }
.sc-tr.danger   { color: #f87171; }
.sc-tr.neutral  { opacity: .7; }

/* ── Chart grids ──────────────────────────────────────────────────────────── */
.g1 { display: grid; grid-template-columns: 1fr;           gap: 10px; margin-bottom: 10px; }
.g2 { display: grid; grid-template-columns: 1fr 1fr;       gap: 10px; margin-bottom: 10px; }
.g3 { display: grid; grid-template-columns: 1fr 1fr 1fr;   gap: 10px; margin-bottom: 10px; }

/* ── Chart card ───────────────────────────────────────────────────────────── */
.card { border-radius: var(--r-lg, 16px); padding: 16px; transition: border-color .15s; }
.card-hd { display: flex; justify-content: space-between; align-items: flex-start; gap: 12px; margin-bottom: 12px; }
.card h3 { font-size: 13px; font-weight: 500; margin-bottom: 2px; }
.card p  { font-size: 11px; }
.cw      { position: relative; height: 175px; }
.cw-sm   { height: 155px; }
.cw-sla  { height: 155px; }

.leg   { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; flex-shrink: 0; }
.leg-b { margin-bottom: 8px; }
.li    { display: flex; align-items: center; gap: 5px; font-size: 11px; }
.ld    { width: 8px; height: 8px; border-radius: 2px; flex-shrink: 0; }

/* ── SLA cards ────────────────────────────────────────────────────────────── */
.sla-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.sla-c    { border-radius: var(--r-lg, 16px); padding: 16px; transition: border-color .15s; }
.sla-hd   { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.sla-ico  { width: 36px; height: 36px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 16px; flex-shrink: 0; }
.sla-lbl  { font-size: 11px; margin-bottom: 3px; }
.sla-val  { font-size: 22px; font-weight: 600; letter-spacing: -.5px; font-variant-numeric: tabular-nums; }
.sla-trk  { height: 4px; border-radius: 2px; overflow: hidden; }
.sla-fill { height: 100%; border-radius: 2px; transition: width .5s ease; }
.sla-ft   { display: flex; justify-content: space-between; font-size: 11px; margin-top: 7px; }

/* ── Spinner & animations ─────────────────────────────────────────────────── */
.spinner { width: 14px; height: 14px; border: 1.5px solid rgba(255,255,255,.15); border-top-color: #a5b4fc; border-radius: 50%; animation: spin .6s linear infinite; }
@keyframes spin  { to { transform: rotate(360deg); } }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: .4; } }
.pulse { animation: pulse 1.8s infinite; }
</style>

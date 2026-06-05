<script setup>
import { computed, onMounted, ref } from "vue";
import {
  Chart as ChartJS,
  ArcElement,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
} from "chart.js";
import { Doughnut, Bar } from "vue-chartjs";

import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

ChartJS.register(
  ArcElement,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend
);

const stats = ref({
  total_assets: 0,
  active_assets: 0,
  total_incidents: 0,
  open_incidents: 0,
  in_progress_incidents: 0,
  resolved_incidents: 0,
  closed_incidents: 0,
  low_incidents: 0,
  medium_incidents: 0,
  high_incidents: 0,
  critical_incidents: 0,
  overdue_incidents: 0,
  unread_notifications: 0,
});

const fetchStats = async () => {
  try {
    const token = localStorage.getItem("token");

    const response = await api.get("/dashboard/stats", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    stats.value = response.data;
  } catch (error) {
    console.error(error);
  }
};

const statusChartData = computed(() => ({
  labels: ["Open", "In Progress", "Resolved", "Closed"],
  datasets: [
    {
      data: [
        stats.value.open_incidents,
        stats.value.in_progress_incidents,
        stats.value.resolved_incidents,
        stats.value.closed_incidents,
      ],
    },
  ],
}));

const priorityChartData = computed(() => ({
  labels: ["Low", "Medium", "High", "Critical"],
  datasets: [
    {
      label: "Incidents",
      data: [
        stats.value.low_incidents,
        stats.value.medium_incidents,
        stats.value.high_incidents,
        stats.value.critical_incidents,
      ],
    },
  ],
}));

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
};

onMounted(fetchStats);
</script>

<template>
  <DashboardLayout>
    <h1 class="text-2xl font-bold mb-6">Dashboard</h1>

    <div class="grid grid-cols-4 gap-6">
      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Total Assets</h3>
        <p class="text-3xl font-bold mt-2">{{ stats.total_assets }}</p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Assets Actifs</h3>
        <p class="text-3xl font-bold mt-2">{{ stats.active_assets }}</p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Total Incidents</h3>
        <p class="text-3xl font-bold mt-2">{{ stats.total_incidents }}</p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Notifications non lues</h3>
        <p class="text-3xl font-bold mt-2 text-blue-600">
          {{ stats.unread_notifications }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Incidents ouverts</h3>
        <p class="text-3xl font-bold mt-2 text-orange-600">
          {{ stats.open_incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">En cours</h3>
        <p class="text-3xl font-bold mt-2 text-indigo-600">
          {{ stats.in_progress_incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Résolus</h3>
        <p class="text-3xl font-bold mt-2 text-green-600">
          {{ stats.resolved_incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Clôturés</h3>
        <p class="text-3xl font-bold mt-2 text-slate-600">
          {{ stats.closed_incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Critiques</h3>
        <p class="text-3xl font-bold mt-2 text-red-600">
          {{ stats.critical_incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">SLA en retard</h3>
        <p class="text-3xl font-bold mt-2 text-red-700">
          {{ stats.overdue_incidents }}
        </p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 mt-8">
      <div class="bg-white rounded-xl shadow p-6 h-96">
        <h2 class="text-lg font-bold mb-4">
          Incidents par statut
        </h2>

        <Doughnut
          :data="statusChartData"
          :options="chartOptions"
        />
      </div>

      <div class="bg-white rounded-xl shadow p-6 h-96">
        <h2 class="text-lg font-bold mb-4">
          Incidents par priorité
        </h2>

        <Bar
          :data="priorityChartData"
          :options="chartOptions"
        />
      </div>
    </div>
  </DashboardLayout>
</template>
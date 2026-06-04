<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const stats = ref({
  total_assets: 0,
  total_incidents: 0,
  open_incidents: 0,
  critical_incidents: 0,
  active_assets: 0,
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

onMounted(fetchStats);
</script>

<template>
  <DashboardLayout>
    <h1 class="text-2xl font-bold mb-6">Dashboard</h1>

    <div class="grid grid-cols-3 gap-6">
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
        <h3 class="text-gray-500">Incidents Ouverts</h3>
        <p class="text-3xl font-bold mt-2">{{ stats.open_incidents }}</p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">Incidents Critiques</h3>
        <p class="text-3xl font-bold mt-2">{{ stats.critical_incidents }}</p>
      </div>
    </div>
  </DashboardLayout>
</template>
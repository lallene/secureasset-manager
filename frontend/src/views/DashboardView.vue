<script setup>
import { onMounted, ref } from "vue";

import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const stats = ref({
  assets: 0,
  incidents: 0,
  users: 0,
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

onMounted(() => {
  fetchStats();
});
</script>

<template>
  <DashboardLayout>
    <div class="grid grid-cols-3 gap-6">
      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">
          Assets
        </h3>

        <p class="text-3xl font-bold mt-2">
          {{ stats.assets }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">
          Incidents
        </h3>

        <p class="text-3xl font-bold mt-2">
          {{ stats.incidents }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow p-6">
        <h3 class="text-gray-500">
          Users
        </h3>

        <p class="text-3xl font-bold mt-2">
          {{ stats.users }}
        </p>
      </div>
    </div>
  </DashboardLayout>
</template>
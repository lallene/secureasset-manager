<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const incidents = ref([]);

const fetchIncidents = async () => {
  const token = localStorage.getItem("token");

  const response = await api.get("/incidents", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  incidents.value = response.data;
};

onMounted(fetchIncidents);
</script>

<template>
  <DashboardLayout>
    <h1 class="text-2xl font-bold mb-6">Incidents</h1>

    <div class="bg-white rounded-xl shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-slate-100">
          <tr>
            <th class="p-3 text-left">Titre</th>
            <th class="p-3 text-left">Type</th>
            <th class="p-3 text-left">Priorité</th>
            <th class="p-3 text-left">Statut</th>
            <th class="p-3 text-left">Asset</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="incident in incidents" :key="incident.ID" class="border-t">
            <td class="p-3">{{ incident.title }}</td>
            <td class="p-3">{{ incident.type }}</td>
            <td class="p-3">{{ incident.priority }}</td>
            <td class="p-3">{{ incident.status }}</td>
            <td class="p-3">{{ incident.asset?.name }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </DashboardLayout>
</template>
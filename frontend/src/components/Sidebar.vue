<script setup>
import { ref, onMounted } from "vue";
import api from "../api/axios";
import { connectWebSocket } from "../services/websocket";

const name = ref("");
const role = ref("");
const unreadNotifications = ref(0);

const fetchNotificationsCount = async () => {
  try {
    const token = localStorage.getItem("token");

    const response = await api.get("/notifications", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    unreadNotifications.value =
      response.data.filter((n) => !n.is_read).length;
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  name.value = localStorage.getItem("name") || "";
  role.value = localStorage.getItem("role") || "";

  fetchNotificationsCount();
  connectWebSocket(() => {
    fetchNotificationsCount();
  });
});

const logout = () => {
  localStorage.clear();
  window.location.href = "/";
};
</script>

<template>
  <aside class="w-64 bg-slate-900 text-white min-h-screen p-5 flex flex-col justify-between">
    <div>
      <h1 class="text-2xl font-bold mb-10">
        SecureAsset
      </h1>

      <nav class="flex flex-col gap-4">

       <router-link
         to="/notifications"
         class="flex items-center justify-between"
       >
         <span>Notifications</span>

         <span
           v-if="unreadNotifications > 0"
           class="bg-red-600 text-white text-xs px-2 py-1 rounded-full"
         >
           {{ unreadNotifications }}
         </span>
       </router-link>

        <router-link to="/dashboard" class="hover:text-blue-400">
          Dashboard
        </router-link>

        <router-link
          v-if="role === 'Admin' || role === 'Technician'"
          to="/assets"
        >
          Assets
        </router-link>

        <router-link to="/incidents">
          Incidents
        </router-link>

        <router-link
          v-if="role === 'Admin'"
          to="/users"
        >
          Users
        </router-link>
      </nav>
    </div>

    <div class="border-t border-slate-700 pt-4">
      <p class="text-sm text-slate-300">
        {{ name }}
      </p>

      <p class="text-xs text-slate-500 mb-3">
        {{ role }}
      </p>

      <button
        @click="logout"
        class="bg-red-600 hover:bg-red-700 text-white px-3 py-2 rounded w-full"
      >
        Déconnexion
      </button>
    </div>
  </aside>
</template>
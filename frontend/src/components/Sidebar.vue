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
      headers: { Authorization: `Bearer ${token}` },
    });
    unreadNotifications.value = response.data.filter((n) => !n.is_read).length;
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
  <aside class="w-64 bg-slate-950 text-slate-200 min-h-screen p-5 flex flex-col justify-between border-r border-slate-900 font-sans antialiased">

    <div>
      <!-- Brand Logo / Identity -->
      <div class="flex items-center gap-2.5 mb-9 px-2">
        <div class="h-8 w-8 flex items-center justify-center rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-500 shadow-md shadow-blue-500/10">
          <svg class="h-4 w-4 text-white" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
          </svg>
        </div>
        <span class="text-lg font-bold tracking-tight text-white">SecureAsset</span>
      </div>

      <!-- Navigation Links -->
      <nav class="space-y-1">

        <!-- Dashboard (Premier de liste, logique et indispensable) -->
        <router-link
          to="/dashboard"
          active-class="bg-slate-900 text-white font-medium"
          class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
          </svg>
          <span class="text-sm">Dashboard</span>
        </router-link>

        <!-- Incidents -->
        <router-link
          to="/incidents"
          active-class="bg-slate-900 text-white font-medium"
          class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
          </svg>
          <span class="text-sm">Incidents</span>
        </router-link>

        <!-- Assets (Conditionné au rôle) -->
        <router-link
            v-if="role === 'Admin' || role === 'Agent'"
          to="/assets"
          active-class="bg-slate-900 text-white font-medium"
          class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-13.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a3 3 0 013-3v0a3 3 0 013 3M16.5 7.5h.008v.008H16.5V7.5z" />
          </svg>
          <span class="text-sm">Assets</span>
        </router-link>

        <!-- Notifications -->
        <router-link
          to="/notifications"
          active-class="bg-slate-900 text-white font-medium"
          class="flex items-center justify-between text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <div class="flex items-center gap-3">
            <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
            </svg>
            <span class="text-sm">Notifications</span>
          </div>

          <span
            v-if="unreadNotifications > 0"
            class="bg-blue-600 text-white text-[11px] font-semibold px-2 py-0.5 rounded-md min-w-5 text-center shadow-sm shadow-blue-600/20"
          >
            {{ unreadNotifications }}
          </span>
        </router-link>

        <div
            v-if="role === 'Admin' || role === 'Agent' || role === 'Requester'"
            class="px-3 pt-4 pb-2 text-[11px] uppercase tracking-wider text-slate-600 font-semibold"
        >
          CMDB
        </div>

        <router-link
            v-if="role === 'Admin' || role === 'Agent' || role === 'Requester'"
            to="/cmdb/applications"
            active-class="bg-slate-900 text-white font-medium"
            class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 7.5h10.5M6.75 12h10.5m-10.5 4.5h10.5M4.5 4.5h15A1.5 1.5 0 0121 6v12a1.5 1.5 0 01-1.5 1.5h-15A1.5 1.5 0 013 18V6a1.5 1.5 0 011.5-1.5z" />
          </svg>
          <span class="text-sm">Applications</span>
        </router-link>
        <router-link
            v-if="role === 'Admin' || role === 'Agent' || role === 'Requester'"
            to="/cmdb/databases"
            active-class="bg-slate-900 text-white font-medium"
            class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg
              class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.75"
              stroke="currentColor"
          >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.75 6.75C3.75 5.507 7.444 4.5 12 4.5s8.25 1.007 8.25 2.25S16.556 9 12 9 3.75 7.993 3.75 6.75z"
            />
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.75 6.75v5.25C3.75 13.243 7.444 14.25 12 14.25s8.25-1.007 8.25-2.25V6.75"
            />
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.75 12v5.25c0 1.243 3.694 2.25 8.25 2.25s8.25-1.007 8.25-2.25V12"
            />
          </svg>

          <span class="text-sm">Bases de données</span>
        </router-link>
        <router-link
            v-if="role === 'Admin' || role === 'Agent' || role === 'Requester'"
            to="/cmdb/relations"
            active-class="bg-slate-900 text-white font-medium"
            class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 7.5h.008v.008H7.5V7.5zm9 9h.008v.008H16.5V16.5zM8.25 8.25l7.5 7.5m.75-8.25a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-9 9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
          </svg>
          <span class="text-sm">Relations</span>
        </router-link>

        <div
            v-if="role === 'Admin'"
            class="px-3 pt-4 pb-2 text-[11px] uppercase tracking-wider text-slate-600 font-semibold"
        >
          Administration
        </div>

        <!-- Users (Conditionné au rôle Admin) -->
        <router-link
          v-if="role === 'Admin'"
          to="/users"
          active-class="bg-slate-900 text-white font-medium"
          class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
          </svg>
          <span class="text-sm">Utilisateurs</span>
        </router-link>
        <router-link
            v-if="role === 'Admin'"
            to="/sites"
            active-class="bg-slate-900 text-white font-medium"
            class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 21h16.5M4.5 3h15a.75.75 0 01.75.75V21H3.75V3.75A.75.75 0 014.5 3zm3 4.5h2.25v2.25H7.5V7.5zm0 4.5h2.25v2.25H7.5V12zm0 4.5h2.25v2.25H7.5V16.5z" />
          </svg>
          <span class="text-sm">Sites</span>
        </router-link>
        <router-link
            v-if="role === 'Admin'"
            to="/services"
            active-class="bg-slate-900 text-white font-medium"
            class="flex items-center gap-3 text-slate-400 hover:text-slate-100 hover:bg-slate-900/50 px-3 py-2.5 rounded-xl transition-all duration-150 group"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.75" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 6.75V4.875c0-.621.504-1.125 1.125-1.125h3.75C14.496 3.75 15 4.254 15 4.875V6.75m-9 0h12m-12 0v10.5A2.25 2.25 0 008.25 19.5h7.5A2.25 2.25 0 0018 17.25V6.75" />
          </svg>
          <span class="text-sm">Services</span>
        </router-link>

      </nav>
    </div>

    <!-- User Profile Footer Section -->
    <div class="border-t border-slate-900 pt-4 flex items-center justify-between px-1">
      <div class="min-w-0 flex-1 pr-2">
        <p class="text-sm font-medium text-white truncate">
          {{ name }}
        </p>
        <p class="text-xs text-slate-500 truncate mt-0.5">
          {{ role }}
        </p>
      </div>

      <!-- Bouton logout discret type "Power off" au lieu d'un gros bloc rouge -->
      <button
        @click="logout"
        title="Déconnexion"
        class="text-slate-500 hover:text-rose-400 hover:bg-rose-500/10 p-2 rounded-xl transition-all duration-150 flex-shrink-0"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
        </svg>
      </button>
    </div>

  </aside>
</template>
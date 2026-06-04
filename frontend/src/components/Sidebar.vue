<template>
  <aside class="w-64 bg-slate-900 text-white min-h-screen p-5">
    <h1 class="text-2xl font-bold mb-10">
      SecureAsset
    </h1>

    <nav class="flex flex-col gap-4">
      <router-link to="/dashboard" class="hover:text-blue-400">Dashboard</router-link>
      <router-link to="/assets" class="hover:text-blue-400">Assets</router-link>
      <router-link to="/incidents" class="hover:text-blue-400">Incidents</router-link>
      <a href="#" class="hover:text-blue-400">
        Users
      </a>
    </nav>
  </aside>
</template>
<script setup>
import Sidebar from "../components/Sidebar.vue";
import Navbar from "../components/Navbar.vue";
</script>

<template>
  <div class="flex bg-slate-100">
    <Sidebar />

    <div class="flex-1">
      <Navbar />

      <main class="p-6">
        <slot />
      </main>
    </div>
  </div>
</template>
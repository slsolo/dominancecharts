<script setup lang="ts">
import {ref, computed, onMounted} from "vue";
const furs = ref(null);
const loading = ref(true);
const error = ref(null);
const first = ref("");
const second = ref("");
const result = ref("");
function fetchCompare() {
  loading.value = true;
  return fetch(`/api/furs/compare/${first.value}/${second.value}`, {
    method: 'get'
  }).then(res => {
    if(!res.ok) {
      const error = new Error(res.statusText);
      throw error;
    }
    return res.text();
  }).then(v => {result.value = v}).catch(err => {
    error.value = err;
  }).then(() => {
    loading.value = false;
  })
}
function fetchFurNames() {
  loading.value = true;
  return fetch("/api/furs/names", {
    method: 'get'
  }).then(res => {
    if(!res.ok) {
      const error = new Error(res.statusText);
      throw error;
    }
    return res.json();
  }).then(json => {furs.value = json}).catch(err => {
    error.value = err;
  }).then(() => {
    loading.value = false;
  })
}
onMounted(() => {
  fetchFurNames();
});
</script>

<template>
  <main>
    <p>first</p>
    <select v-model="first" v-if="!loading && furs">
      <option v-for="(item, key) in furs" :key="key" :value="item">{{item}}</option>
    </select>
    <p>second</p>
    <select v-model="second" v-if="!loading && furs">
      <option v-for="(item, key) in furs" :key="key" :value="item">{{item}}</option>
    </select>
    <button @click="fetchCompare">Check</button>
    <p v-if="!loading && result">{{result}}</p>
    <p v-else>{{error}}</p>
  </main>
</template>

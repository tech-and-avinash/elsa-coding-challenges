<template>
  <div>
    <div class="toolbar">
      <div>
        <h1>Settings</h1>
        <div class="sub">Workspace preferences and defaults.</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-body">
        <div class="field">
          <label>Workspace name</label>
          <input v-model="workspaceName" />
        </div>
        <div class="field">
          <label>Default scoring</label>
          <select v-model="defaultScoring">
            <option>Correctness + response speed</option>
            <option>Correctness only</option>
          </select>
        </div>
        <button class="btn primary" @click="saveSettings">Save settings</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const workspaceName = ref('ELSA Live Quiz Workspace');
const defaultScoring = ref('Correctness + response speed');

const saveSettings = () => {
  alert('Settings saved successfully!');
};
</script>

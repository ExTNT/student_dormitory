<template>
  <main class="forbidden">
    <el-result icon="warning" title="403" sub-title="无权限访问该页面">
      <template #extra>
        <el-button type="primary" @click="$router.back()">返回</el-button>
      </template>
    </el-result>
  </main>
</template>

<style scoped>
.forbidden {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at 50% 30%, rgba(11, 93, 102, 0.13), transparent 28%),
    var(--bg);
}
</style>

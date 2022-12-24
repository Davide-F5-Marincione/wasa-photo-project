<script>
import TopBar from '../components/TopBar.vue';

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
    components: { TopBar }
}
</script>

<template>
	<TopBar></TopBar>
	<div>
        <b>{{ $route.params }}</b>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

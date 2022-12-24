<script>
import PostsHolder from '../components/PostsHolder.vue';

export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
        };
    },
    components: { PostsHolder }
}
</script>

<template>
	<TopBar></TopBar>
	<div>
		<PostsHolder></PostsHolder>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

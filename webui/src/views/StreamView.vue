<script>
export default {
    data: function () {
        return {
            errormsg: null,
            results: new Set(),
            limit: ""
        };
    },
    methods: {
        async refreshData() {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/users/" + localStorage.username + "/stream", {params: {"photos-limit":this.limit}});
                
                response.data.forEach(element => {
                    this.results.add(element)
                });
                if (response.data.length > 0) {
                    this.limit = response.data[response.data.length - 1];
                }

			} catch (e) {
				this.errormsg = e.toString();

			}
		},
        async delPost(id) {
			this.errormsg = null;
			try {
				await this.$axios.delete("/photos/" + id.toString());
                this.results.clear()
                this.limit = ""
                this.refreshData()
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
    },
    created() {
        this.refreshData();
    }
}
</script>

<template>
	<TopBar></TopBar>
	<div>
		<div class="posts-holder">
            <PostCard v-for="elem in results" :imgId="elem" :del="()=>delPost(elem)" v-bind:key="elem"></PostCard>
            <button class="posts-more" v-on:click="refreshData">Show more posts!</button>
        </div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>
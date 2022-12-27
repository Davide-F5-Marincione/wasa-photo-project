<script>
export default {
    data: function () {
        return {
            errormsg: null,
            results: [],
            limit: "",
            resp: ""
        };
    },
    methods: {
        async refreshData() {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/users/" + localStorage.username + "/stream", {params: {"photos-limit":this.limit}});
                
                this.resp = response.data

                this.results.push(...response.data)
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
                this.results = []
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
            <PostCard v-for="elem in results" v-bind:imgId="elem" v-bind:del="()=>delPost(elem)"></PostCard>
            <button class="posts-more" :onclick="() => refreshData()">Show more posts!</button>
        </div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>
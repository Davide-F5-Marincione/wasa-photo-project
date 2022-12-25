<script>
import char from "../assets/default_char.jpg"
export default {
	props: ["title", "imgId"],
	data: function () {
        return {
            errormsg: null,
			src: char,
			resp: ""
        };
    },
	methods: {
		async getSource() {
			this.errormsg = null;
			const urlCreator = window.URL || window.webkitURL;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString() + "/raw", {responseType: 'arraybuffer'});
				var img = Buffer.from(response.data, 'binary').toString('base64');
				this.src = "data:" + response.headers["content-type"] + ";base64, " + img
			} catch (e) {
				this.errormsg = e.toString();
			}
		}
	},
	created() {
		if (this.imgId) {
			this.getSource();
		}
	}
}
</script>


<template>
	<div class="post-card">
		<div class="post-title">
			{{ title }}
		</div>
		<div class="post-photoandcomments">
			<div class="post-photo-container">
				<img class="post-photo" :alt="title" :src="src"/>
			</div>
			<div class="post-comments">
				{{ resp }}
			</div>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>
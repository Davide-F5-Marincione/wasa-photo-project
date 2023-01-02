<script>
import travolta from "../assets/confused_travolta.gif" // Is there any reason for defaulting to this? No.
export default {
	props: ["alt", "imgId"],
	data: function () {
        return {
            errormsg: null,
			src: "#",
			resp: ""
        };
    },
	methods: {
		async getSource() {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString() + "/raw", {responseType: 'arraybuffer'});

				// Got the idea from here https://stackoverflow.com/questions/8499633/how-to-display-base64-images-in-html
				var img = Buffer.from(response.data, 'binary').toString('base64');
				this.src = "data:" + response.headers["content-type"] + ";base64, " + img;
			} catch (e) {
				this.src = travolta;
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
	<div class="post-photo-container">
        <img class="post-photo" :alt="alt" :src="src"/>
    </div>
</template>
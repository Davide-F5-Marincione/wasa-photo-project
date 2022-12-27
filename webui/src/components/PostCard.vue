<script>
import ImageContainer from "./ImageContainer.vue";
export default {
    props: ["imgId"],
    data: function () {
        return {
            errormsg: null,
			title:"Wait... is something missing?",
			author:"The browser itself",
			date: new Date().toLocaleString( 'sv', { timeZoneName: 'short' } ).split(" ").slice(0,2).join(" "),
			liked: false,
			likersResult: [],
			commentsResults: [],
			likesBase: "",
			commentsLimit: "",
			username: localStorage.username
        };
    },
    methods: {
		async refreshData(new_likes=false, new_comments=false) {
			if (!this.imgId) {
				return
			}
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString(), {params: {"likes-base": this.likesBase, "comments-limit":this.commentsLimit}});
				this.title = response.data["photo-title"]
				this.author = response.data["photo-author"]
				this.date = this.convDate(response.data["photo-date"])
				this.liked = response.data["liked"]
				
				for (let i=0; i< response.data["likes-running-batch"].length; i++) {
					if (response.data["likes-running-batch"][i].name == this.username) {
						response.data["likes-running-batch"].splice(i, 1)
						break
					}
				}

				if (new_likes) {
					this.likersResult.push(...response.data["likes-running-batch"])
					if (response.data["likes-running-batch"].length > 0) {
						this.likesBase = response.data["likes-running-batch"][response.data["likes-running-batch"].length - 1].name;
					}
				}
				if (new_comments) {
					this.commentsResults.push(...response.data["comments-running-batch"])
					if (response.data["comments-running-batch"].length > 0) {
						this.commentsLimit = response.data["comments-running-batch"][response.data["comments-running-batch"].length - 1]["comment-id"];
					}
				}

			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async likeThis() {
			if (!this.imgId) {
				return
			}
			try {
				await this.$axios.put("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
		},
		async unlikeThis() {
			if (!this.imgId) {
				return
			}
			try {
				await this.$axios.delete("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
		},
		convDate(date) {
			return date //TODO: convert to local timezone
		},
    },
	created() {
		if (this.imgId) {
			this.refreshData(true, true);
		}
	},
    components: { ImageContainer }
}
</script>


<template>
	<div class="post-card">
		<div class="post-titleinfoandphoto">
			<div class="post-titleandinfo">
				<div class="post-title">
					{{ title }}
				</div>
				<div class="post-info">
					<div>					
						<router-link class="post-author" align="right" :to="'/users/' + author">{{ author }}</router-link>
						<div class="post-date" align="right">
							{{ date }}
						</div>
					</div>
				</div>
			</div>
			<ImageContainer :alt="title" :img-id="imgId"></ImageContainer>
		</div>
		<div class="post-commentsandlikes">
			<div id="likesholder" class="post-likesholder">
				<router-link class="like-element" v-for="element in likersResult" :to="'/users/' +element['name']">{{ element['name'] }}</router-link>
				<router-link v-if="liked" class="like-element" :to="'/users/' + username">{{ username }}</router-link>
				<button class="like-element-end" :onclick="() => this.refreshData(new_likes=true)">More likers!</button>
			</div>
			<button v-if="!liked" class="post-like-button" :onclick="() => this.likeThis().then()">like</button>
			<button v-if="liked" class="post-unlike-button" :onclick="() => this.unlikeThis().then()">unlike</button>
			<div class="post-comments">
			</div>
			<textarea class="post-comment-writer" maxlength="256"></textarea>
			<button id="button-comment" class="post-comment-button">
				Comment!
			</button> 
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>
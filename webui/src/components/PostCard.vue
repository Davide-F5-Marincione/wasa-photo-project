<script>
import ImageContainer from "./ImageContainer.vue";
export default {
    props: ["imgId", "del"],
    data: function () {
        return {
            errormsg: null,
			title:"Wait... is something missing?",
			author:"The browser itself",
			date: new Date().toLocaleString( 'sv', { timeZoneName: 'short' } ).split(" ").slice(0,2).join(" "),
			liked: false,
			likersResult: new Set(),
			commentsResults: [],
			likesBase: "",
			commentsLimit: "",
			username: localStorage.username,
			commentWrite: ""
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

				if (new_likes) {
					response.data["likes-running-batch"].forEach(element => {
						if (element.name != this.username) {
							this.likersResult.add(element.name)
						}
					});
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
			return (new Date(date.split(" ").join("T") + "Z")).toLocaleString( 'sv', { timeZoneName: 'short' } ).split(" ").slice(0,2).join(" ") //TODO: convert to local timezone
		},
		async sendComment() {
			if (!this.imgId) {
				return
			}
			if (this.commentWrite.length <= 0) {
				return
			}
			try {
				await this.$axios.post("/photos/" + this.imgId.toString() + "/comments",  this.commentWrite, {headers: {'Content-Type': 'application/json'}});
				this.commentWrite = ""
				this.commentsResults = []
				this.commentsLimit = ""
				this.refreshData(false, true)
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		async deleteComment(id) {
			if (!this.imgId) {
				return
			}

			try {
				await this.$axios.delete("/photos/" + this.imgId.toString() + "/comments/" + id.toString());
				this.commentsResults = []
				this.commentsLimit = ""
				this.refreshData(false, true)
			} catch (e) {
				this.errormsg = e.toString();
			}
		}
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
						<button v-if="author != username && !liked" class="post-like-button" v-on:click="likeThis">like</button>
						<button v-if="author != username && liked" class="post-unlike-button" v-on:click="unlikeThis">unlike</button>
						<button v-if="author == username" class="delete-post" align="right" v-on:click="del">delete this post</button>
					</div>
				</div>
			</div>
			<ImageContainer :alt="title" :img-id="imgId"></ImageContainer>
		</div>
		<div class="post-commentsandlikes">
			<div id="likesholder" class="post-likesholder disable-scrollbars">
				Users who like this:
				<router-link class="like-element" v-for="element in likersResult" v-bind:key="element" :to="'/users/' +element">{{ element }}</router-link>
				<router-link v-if="liked" class="like-element" :to="'/users/' + username">{{ username }}</router-link>
				<button class="like-element-end" v-on:click="() => this.refreshData(true, false)">More likes!</button>
			</div>
			<div class="post-comments disable-scrollbars">
				Users' comments:
				<div v-for="element in commentsResults" class="post-comment" v-bind:key="element['comment-id']">
					<div class="comment-authoranddeleter">
						<router-link class="comment-author" :to="'/users/' + element['comment-author']"> {{ element["comment-author"] }}: </router-link>
						<button v-if="element['comment-author']==username" class="comment-deleter" v-on:click="() => this.deleteComment(element['comment-id'])">✕</button>
					</div>
					<div class="comment-content">{{ element["comment-text"] }}</div>
					<div class="comment-date">{{ convDate(element["comment-date"]) }}</div>
				</div>
				<button class="post-comments-end" v-on:click="() => this.refreshData(false, true)">See more comments!</button>
			</div>
			<textarea class="post-comment-writer" maxlength="256" v-model="commentWrite"></textarea>
			<button class="post-comment-button" v-on:click="sendComment">
				Comment!
			</button>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>
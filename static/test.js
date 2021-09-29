var app = new Vue({
  el: "#app",
  data() {
    return {
      listVariants: {},
      showList: true,
      is_finished: false,
      current_variant_id: "",
      current_task_id: "",
      current_test_id: "",
      current_question: "",
      current_answers: [],
      current_percent: "",
    };
  },
  mounted() {
    this.getListVariants();
  },
  methods: {
    toStart() {
      this.getListVariants();
      this.is_finished = false;
      this.showList = true;
    },
    singout() {
      axios
        .put("/api")
        .then((response) => {
          window.location.href = "/";
        })
        .catch((error) => {
          console.log(error);
          alert(error.response.data);
        });
    },
    getListVariants() {
      axios
        .get("/api")
        .then((response) => {
          this.listVariants = response.data.variants;
        })
        .catch((error) => {
          console.log(error);
          alert(error.response.data);
        });
    },
    sendAnswer(id) {
      params = {
        variant_id: this.current_variant_id,
        task_id: this.current_task_id,
        test_id: this.current_test_id,
        answer: Number(id),
      };
      console.log(params);
      axios
        .post("/api", params)
        .then((response) => {
          if (response.data.status === "ok") {
            this.get_task_by_id(
              this.current_variant_id,
              this.current_task_id + 1
            );
          } else if (response.data.status === "finished") {
            this.current_percent = response.data.percent;
            this.is_finished = true;
          }
        })
        .catch((error) => {
          console.log(error);
          alert(error.response.data);
          if (error.response.data.includes("question already answered")) {
            this.get_task_by_id(
              this.current_variant_id,
              this.current_task_id + 1
            );
          }
        });
    },
    get_task_by_id(variant_id, task_id) {
      axios
        .get("/api", {
          params: {
            variant_id: variant_id,
            task_id: task_id,
          },
        })
        .then((response) => {
          console.log(response.data.question);
          this.current_task_id = response.data.question.id;
          this.current_variant_id = response.data.question.variant_id;
          this.current_test_id = response.data.question.test_id;
          this.current_question = response.data.question.question;
          this.current_answers = response.data.question.answers;
          this.showList = false;
        })
        .catch(function (error) {
          console.log(error);
          alert(error.response.data);
        });
    },
  },
});

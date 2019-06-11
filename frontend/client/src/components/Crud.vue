<template>
  <v-container>
    <v-flex xs12 sm4 md4 offset-xs4>
      <v-text-field
          label="端到端密钥"
          hint="密码仅保存在浏览器端"
          v-model="e2eePass"
      ></v-text-field>
    </v-flex>

    <div id="app">
      <v-app id="demo">
        <div>
          <v-toolbar flat color="white">
            <v-toolbar-title>疫苗供应链管理 Demo</v-toolbar-title>
            <v-divider
                class="mx-2"
                inset
                vertical
            ></v-divider>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialog" max-width="500px">
              <template v-slot:activator="{ on }">
                <v-btn color="#666" dark class="mb-2" v-on="on">Add</v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">{{ formTitle }}</span>
                </v-card-title>

                <v-flex xs12 sm6 md4>

                </v-flex>

                <v-card-text>
                  <v-container grid-list-md>
                    <v-layout wrap>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.serial" label="Serial"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.factory" label="Factory"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-menu
                            ref="menu"
                            v-model="menu"
                            :close-on-content-click="false"
                            :nudge-right="40"
                            :return-value.sync="date"
                            lazy
                            transition="scale-transition"
                            offset-y
                            full-width
                            min-width="290px"
                        >
                          <template v-slot:activator="{ on }">
                            <v-text-field
                                v-model="editedItem.date"
                                label="Date"
                                readonly
                                v-on="on"
                            ></v-text-field>
                          </template>
                          <v-date-picker v-model="editedItem.date" no-title scrollable>
                            <v-spacer></v-spacer>
                            <v-btn flat color="primary" @click="menu = false">Cancel</v-btn>
                            <v-btn flat color="primary" @click="$refs.menu.save(date)">OK</v-btn>
                          </v-date-picker>
                        </v-menu>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.batch" label="Batch"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.carrier" label="Carrier"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.cold_van" label="ColdVan"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.distributor" label="Distributor"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.hospital" label="Hospital"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.patient" label="Patient"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field disabled v-model="editedItem.attach_uri" label="Attach"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md8>
                        <v-text-field disabled v-model="editedItem.attach_sum" label="CheckSum"></v-text-field>
                      </v-flex>
                    </v-layout>
                  </v-container>
                </v-card-text>
                <template>
                  <div id="dropUpload">
                    <vue-dropzone id="uploadField" ref="myVueDropzone"
                                  @vdropzone-success="UploadSuccess"
                                  thumbnailHeight="80"
                                  thumbnailWidth="80"
                                  :options="dropOptions">
                    </vue-dropzone>
                  </div>
                </template>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" flat @click="close">Cancel</v-btn>
                  <v-btn color="blue darken-1" flat @click="save">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-toolbar>
          <v-data-table
              :headers="headers"
              :items="items"
              :expand="expand"
              :loading="loadingMain"
              :rows-per-page-items="rowsPerPage"
              class="crud"
          >
            <v-progress-linear v-slot:progress color="blue" indeterminate></v-progress-linear>
            <template slot="items" slot-scope="props">
              <tr @click="expendItem(props)">
                <td>{{ props.item.id }}</td>
                <td class="text-xs-left">{{ props.item.serial }}</td>
                <td class="text-xs-left factory-cell">{{ props.item.factory }}</td>
                <td class="text-xs-left">{{ props.item.date }}</td>
                <td class="text-xs-left">{{ props.item.batch }}</td>
                <td class="text-xs-left">{{ props.item.carrier }}</td>
                <td class="text-xs-left">{{ props.item.cold_van }}</td>
                <td class="text-xs-left">{{ props.item.distributor }}</td>
                <td class="text-xs-left">{{ props.item.hospital }}</td>
                <td class="text-xs-left patient-cell">{{ props.item.patient }}</td>
                <td class="justify-center px-0">
                  <span v-if="props.item.attach_uri"><v-icon color="#bbb" small>attach_file</v-icon></span>
                </td>
                <td class="justify-center layout px-0">
                  <v-icon
                      small
                      class="mr-1"
                      @click="editItem(props.item)"
                  >
                    edit
                  </v-icon>
                  <v-icon
                      small
                      class="mr-3"
                      @click="deleteItem(props.item)"
                  >
                    delete
                  </v-icon>
                </td>
              </tr>
            </template>
            <template v-slot:expand="props" v-slot:no-data>
              <v-layout row>
                <v-flex xs12 sm6 md9 offset-sm1>
                  <v-list
                      :loading="false"
                      dense
                      two-line
                  >
                    <template v-for="(item, index) in expandItems">
                      <v-list-tile
                          :key="item.id"
                          avatar
                          ripple
                          @click="goExplorer(item.hash)"
                      >
                        <v-list-tile-content>
                          <v-icon small>verified_user</v-icon>
                          <v-list-tile-title>
                            <div>
                              <prism language="sql" :code="item.sql"></prism>
                            </div>
                          </v-list-tile-title>

                          <v-list-tile-sub-title>
                            <v-icon small>fingerprint</v-icon>
                            {{ item.hash }}
                          </v-list-tile-sub-title>
                        </v-list-tile-content>

                        <v-list-tile-action>
                          <v-list-tile-action-text>{{ item.user }}</v-list-tile-action-text>
                          <v-icon>
                            storage
                          </v-icon>
                        </v-list-tile-action>

                      </v-list-tile>
                      <v-divider
                          v-if="index + 1 < items.length"
                          :key="index"
                      ></v-divider>
                    </template>
                  </v-list>
                </v-flex>
              </v-layout>
            </template>
            <template v-slot:no-data>
              <v-btn color="primary" @click="initialize">Reload</v-btn>
            </template>
          </v-data-table>
        </div>
      </v-app>
    </div>
  </v-container>
</template>

<script>
  import axios from 'axios'
  import vueDropZone from "vue2-dropzone"
  import 'prismjs'
  import 'prismjs/themes/prism.css'
  import Prism from 'vue-prism-component'
  import 'prismjs/components/prism-sql'
  import e2e from 'e2e_js'

  const aes = require('aes-js');

  export default {
    data: () => ({
      e2eePass: '',
      dropOptions: {
        url: "/apiv1/attach",
        maxFilesize: 20, // MB
        maxFiles: 1,
        chunking: false,
        chunkSize: 500, // Bytes
        thumbnailWidth: 100, // px
        thumbnailHeight: 100,
        addRemoveLinks: true,
        headers: {"Authorization": localStorage.getItem('token')}
      },
      haveAttach: false,
      isMounted: false,
      expand: false,
      loadingMain: true,
      rowsPerPage: [20, 50, 100, {"text": "$vuetify.dataIterator.rowsPerPageAll", "value": -1}],
      menu: '',
      date: '',
      picker: '',
      formatted: '',
      dialog: false,
      headers: [
        {text: 'ID', value: 'id', align: 'left'},
        {text: 'Serial', value: 'serial', align: 'left'},
        {text: 'Factory', value: 'factory', align: 'left'},
        {text: 'Date', value: 'date'},
        {text: 'Batch', value: 'batch'},
        {text: 'Carrier', value: 'carrier'},
        {text: 'ColdVan', value: 'cold_van'},
        {text: 'Distributor', value: 'distributor'},
        {text: 'Hospital', value: 'hospital'},
        {text: 'Patient', value: 'patient'},
      ],
      items: [],
      expandItems: [],
      editedIndex: -1,
      editedItem: {
        serial: '',
        factory: '',
        date: '',
        batch: '',
        carrier: '',
        cold_van: '',
        distributor: '',
        hospital: '',
        patient: '',
        attach_uri: '',
        attach_sum: '',
      },
      defaultItem: {
        serial: '',
        factory: '',
        date: '',
        batch: '',
        carrier: '',
        cold_van: '',
        distributor: '',
        hospital: '',
        patient: '',
        attach_uri: '',
        attach_sum: '',
      }
    }),

    components: {
      vueDropzone: vueDropZone,
      Prism
    },

    computed: {
      formTitle() {
        return this.editedIndex === -1 ? 'New Item' : 'Edit Item'
      }
    },
    mounted() {
      if (localStorage.e2eePass) {
        this.e2eePass = localStorage.e2eePass;
      }
    },
    watch: {
      dialog(val) {
        val || this.close()
      },
      e2eePass(pass) {
        localStorage.e2eePass = pass;
      }
    },

    created() {
      this.initialize()
    },

    methods: {
      to_hex(d) {
        return aes.utils.hex.fromBytes(d);
      },
      from_hex(s) {
        return new Uint8Array(aes.utils.hex.toBytes(s));
      },
      UploadZoneMount() {
        let file = {size: 0, name: this.editedItem.attach_uri};
        let url = "/apiv1/attach/" + this.editedItem.attach_uri;
        this.$log.debug(file, url)
        this.$refs.myVueDropzone.manuallyAddFile(file, url)
        this.isMounted = true;
      },
      UploadSuccess(file, response) {
        this.success = true;
        this.$log.debug(response);
        this.editedItem.attach_uri = response.attach_uri;
        this.editedItem.attach_sum = response.attach_sum;
      },
      initialize() {
        this.loadingMain = true;
        axios.defaults.headers.common['Authorization'] = localStorage.getItem('token');
        axios.get('/apiv1/cargo')
          .then(resp => {
            this.items = resp.data;
            if (resp.data == null) {
              this.$notify({
                group: 'crud',
                type: 'error',
                title: "获取信息失败",
                text: "no cargo got",
              });
            }
            this.loadingMain = false;
            this.$log.debug(resp);
          })
          .catch(err => {
            this.$notify({
              group: 'crud',
              type: 'error',
              title: "获取信息失败",
              text: err.response.data,
            });
            this.loadingMain = false;
            this.$log.error(err);
          });
      },
      expendItem(exp) {
        this.loadingMain = true;
        axios.defaults.headers.common['Authorization'] = localStorage.getItem('token');
        axios.get('/apiv1/sql', {params: {cargo_id: exp.item.id, limit: 10}})
          .then(resp => {
            this.expandItems = resp.data;
            if (resp.data !== null) {
              exp.expanded = !exp.expanded;
            }
            this.loadingMain = false;
            this.$log.debug(resp);
          })
          .catch(err => {
            this.$notify({
              group: 'crud',
              type: 'error',
              title: "获取区块链上的 SQL 记录失败",
              text: err.response.data,
            });
            this.loadingMain = false;
            this.$log.error(err);
          });
      },
      goExplorer(hash) {
        const explorer = 'http://localhost:8082/dbs/' +
          'b77422b30688fdc8facfe84a0c48c1f94aca3444178a9502753b3692a5576f10/requests/';
        window.open(explorer + hash, "_blank");
      },
      editItem(item) {
        axios.defaults.headers.common['Authorization'] = localStorage.getItem('token');
        this.editedIndex = this.items.indexOf(item);
        // this.editedItem = Object.assign({}, item);
        for (let key in item) {
          let value = item[key];
          if (key === 'patient' && value !== null) {
            try {
              this.editedItem[key] = e2e.decrypt_string(this.from_hex(value), this.e2eePass)
            } catch (e) {
              this.$notify({
                group: 'crud',
                type: 'error',
                title: "解密敏感信息失败",
                text: e,
              });
            }
          } else {
            this.editedItem[key] = value
          }
        }

        this.$log.debug("edit: ", this.editedIndex, item);
        this.dialog = true;
        if (item['attach_uri'] !== null) {
          this.UploadZoneMount();
        }
      },

      deleteItem(item) {
        const index = this.items.indexOf(item);
        const id = item.id;
        confirm('Are you sure you want to delete this item?') &&
        axios.delete('/apiv1/cargo/' + id)
          .then(resp => {
            this.$log.debug(resp)
            this.items.splice(index, 1);
          })
          .catch(err => {
            this.$notify({
              group: 'crud',
              type: 'error',
              title: "删除信息失败",
              text: err.response.data,
            });
            this.$log.error(err)
          });
      },

      close() {
        this.dialog = false;
        setTimeout(() => {
          this.editedItem = Object.assign({}, this.defaultItem);
          this.$refs.myVueDropzone.removeAllFiles();
          this.editedIndex = -1;
        }, 300)
      },

      save() {
        if (this.editedIndex > -1) {
          let diff = {};
          for (let key in this.editedItem) {
            let originItem = this.items[this.editedIndex];
            if (!(key in originItem) ||
              this.editedItem[key] !== originItem[key]) {
              let value = this.editedItem[key];
              this.$log.debug(value, originItem[key]);
              if (key === 'patient') {
                if (originItem[key] !== null) {
                  let originValue = e2e.decrypt_string(
                    this.from_hex(originItem[key]),
                    this.e2eePass
                  );
                  if (value === originValue) {
                    continue
                  }
                }
                let encodedVal = new TextEncoder().encode(value);
                diff[key] = this.to_hex(e2e.encrypt(encodedVal, this.e2eePass))
              } else {
                diff[key] = value;
              }
            }
          }
          if (Object.keys(diff).length > 0) {
            diff['id'] = this.editedItem.id;
            this.$log.debug(diff);
            axios.put('/apiv1/cargo', diff)
              .then(resp => {
                this.$log.debug(resp);
                this.initialize();
                // Object.assign(this.items[this.editedIndex], this.editedItem);
                this.$refs.myVueDropzone.removeAllFiles();
              })
              .catch(err => {
                this.$notify({
                  group: 'crud',
                  type: 'error',
                  title: "修改信息失败",
                  text: err.response.data,
                });
                this.$log.error(err)
              });
          }
        } else {
          let insert = {};
          for (let key in this.editedItem) {
            if (key in this.defaultItem &&
              this.editedItem[key] !== this.defaultItem[key]) {
              let value = this.editedItem[key];
              if (key === 'patient') {
                let encodedVal = new TextEncoder().encode(value);
                insert[key] = this.to_hex(e2e.encrypt(encodedVal, this.e2eePass))
              } else {
                insert[key] = value;
              }
            }
          }
          if (Object.keys(insert).length > 0) {
            this.$log.debug(insert)
            axios.post('/apiv1/cargo', insert)
              .then(resp => {
                this.$log.debug(resp)
                this.items.push(this.editedItem)
                this.initialize();
              })
              .catch(err => {
                this.$notify({
                  group: 'crud',
                  type: 'error',
                  title: "创建新条目失败",
                  text: err.response.data,
                });
                this.$log.error(err)
              });
          }
        }
        this.close()
      }
    }
  }
</script>

<style>
  code[class="language-sql"], pre[class="language-sql"] {
    padding: 0;
    margin: 0 0 0 .5em;
    overflow: auto;
    background: hsla(0, 0%, 100%, 0);
  }

  td.patient-cell {
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    max-width: 10em;
  }

  td.factory-cell {
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    max-width: 3em;
  }

  div.dz-details, div.dz-preview {
    width: 120px;
    height: 120px;
  }
</style>

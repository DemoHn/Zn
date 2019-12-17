+ api
  * label: 创建预算
  * method: GET
  * url: activity/create-budget  
  * sig: WEB, CLI
  
  > 须部署在安全域内

  + request    
    > pb/name: `CreateBudgetReq`

    > pb/diff: true

    > 

    - is_online: false (bool?) - 是否已上线
    - status: 1 (num@BudgetStatus) - 预算状态
    - ext_info: (obj@ExtInfo?) - 扩展字段

  + response
    * ref: Budget

+ api
  * label: 使用预算

+ model: ExtInfo
  + oneof
    * desc: 当奖品为「股票卡」时

    - app_key: abcd1234567890 (string) - 股票卡app_key
    - scene_key: abcd9876543210 (string) - 股票卡scene_key

  + oneof
    * desc: 

## section/
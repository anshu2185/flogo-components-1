import { NgModule } from "@angular/core";
import { TcmPubActivityUIContributionHandler } from "./activity";
import { WiContributionHandler } from "wi-studio/app/contrib/wi-contrib";
@NgModule(
  {
    providers: [
      {
        provide: WiContributionHandler,
        useClass: SampleActivityUIContributionHandler
      }
    ]
  }
)
export default class SampleActivityModule {

}

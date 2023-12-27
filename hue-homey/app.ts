import Homey from 'homey';
import { getSceneIdByName, activateScene, getSmartSceneIdByName, activateSmartScene } from './managers/hue-api-manager';

class MyApp extends Homey.App {

  async onInit() {
    const activateSceneCard = this.homey.flow.getActionCard("activate-scene");

    activateSceneCard.registerRunListener(async (args, state) => {
      try {
        const sceneId = await getSceneIdByName(args.scene);
        return activateScene(sceneId);
      } catch (error) {
        console.log(error);
        this.error(error);
      }
    });

    const activateSmartSceneCard = this.homey.flow.getActionCard("activate-smart-scene");

    activateSmartSceneCard.registerRunListener(async (args, state) => {
      try {
        const sceneId = await getSmartSceneIdByName(args.scene);
        return activateSmartScene(sceneId);
      } catch (error) {
        console.log(error);
        this.error(error);
      }
    });
  }
}

module.exports = MyApp;

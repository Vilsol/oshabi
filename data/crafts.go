package data

import (
	"encoding/json"
	"github.com/masatana/go-textdistance"
	"strings"
)

type HarvestType string

//goland:noinspection GoUnusedConst
const (
	// Reveals a random Socket colour crafting effect when Harvested

	HarvestReforgeNonRedToRed        = HarvestType("ReforgeNonRedToRed")
	HarvestReforgeNonBlueToBlue      = HarvestType("ReforgeNonBlueToBlue")
	HarvestReforgeNonGreenToGreen    = HarvestType("ReforgeNonGreenToGreen")
	HarvestReforgeTwoRandomRedBlue   = HarvestType("ReforgeTwoRandomRedBlue")
	HarvestReforgeTwoRandomRedGreen  = HarvestType("ReforgeTwoRandomRedGreen")
	HarvestReforgeTwoRandomBlueGreen = HarvestType("ReforgeTwoRandomBlueGreen")
	HarvestReforgeThreeRandomRGB     = HarvestType("ReforgeThreeRandomRGB")
	HarvestReforgeWhite              = HarvestType("ReforgeWhite")

	// Reforge a Normal, Magic or Rare item as a Rare item with random modifiers, including a * modifier

	HarvestReforgeCaster    = HarvestType("ReforgeCaster")
	HarvestReforgePhysical  = HarvestType("ReforgePhysical")
	HarvestReforgeFire      = HarvestType("ReforgeFire")
	HarvestReforgeAttack    = HarvestType("ReforgeAttack")
	HarvestReforgeLife      = HarvestType("ReforgeLife")
	HarvestReforgeCold      = HarvestType("ReforgeCold")
	HarvestReforgeSpeed     = HarvestType("ReforgeSpeed")
	HarvestReforgeDefence   = HarvestType("ReforgeDefence")
	HarvestReforgeLightning = HarvestType("ReforgeLightning")
	HarvestReforgeChaos     = HarvestType("ReforgeChaos")
	HarvestReforgeCritical  = HarvestType("ReforgeCritical")
	HarvestReforgeInfluence = HarvestType("ReforgeInfluence")

	// Reforge a Normal, Magic or Rare item as a Rare item with random modifiers, including a * modifier. * modifiers are more common

	HarvestReforgeCasterMoreLikely    = HarvestType("ReforgeCasterMoreLikely")
	HarvestReforgePhysicalMoreLikely  = HarvestType("ReforgePhysicalMoreLikely")
	HarvestReforgeFireMoreLikely      = HarvestType("ReforgeFireMoreLikely")
	HarvestReforgeAttackMoreLikely    = HarvestType("ReforgeAttackMoreLikely")
	HarvestReforgeLifeMoreLikely      = HarvestType("ReforgeLifeMoreLikely")
	HarvestReforgeColdMoreLikely      = HarvestType("ReforgeColdMoreLikely")
	HarvestReforgeSpeedMoreLikely     = HarvestType("ReforgeSpeedMoreLikely")
	HarvestReforgeDefenceMoreLikely   = HarvestType("ReforgeDefenceMoreLikely")
	HarvestReforgeLightningMoreLikely = HarvestType("ReforgeLightningMoreLikely")
	HarvestReforgeChaosMoreLikely     = HarvestType("ReforgeChaosMoreLikely")
	HarvestReforgeCriticalMoreLikely  = HarvestType("ReforgeCriticalMoreLikely")
	HarvestReforgeInfluenceMoreLikely = HarvestType("ReforgeInfluenceMoreLikely")

	// Remove a random non-* modifier from a * item and add a new * modifier

	HarvestRemoveAddNonCaster    = HarvestType("RemoveAddNonCaster")
	HarvestRemoveAddNonPhysical  = HarvestType("RemoveAddNonPhysical")
	HarvestRemoveAddNonFire      = HarvestType("RemoveAddNonFire")
	HarvestRemoveAddNonAttack    = HarvestType("RemoveAddNonAttack")
	HarvestRemoveAddNonLife      = HarvestType("RemoveAddNonLife")
	HarvestRemoveAddNonCold      = HarvestType("RemoveAddNonCold")
	HarvestRemoveAddNonSpeed     = HarvestType("RemoveAddNonSpeed")
	HarvestRemoveAddNonDefence   = HarvestType("RemoveAddNonDefence")
	HarvestRemoveAddNonLightning = HarvestType("RemoveAddNonLightning")
	HarvestRemoveAddNonChaos     = HarvestType("RemoveAddNonChaos")
	HarvestRemoveAddNonCritical  = HarvestType("RemoveAddNonCritical")
	HarvestRemoveAddNonInfluence = HarvestType("RemoveAddNonInfluence")
	HarvestRemoveAddInfluence    = HarvestType("RemoveAddInfluence")

	// Augment a * Magic or Rare item with a new * modifier

	HarvestAugmentCaster    = HarvestType("AugmentCaster")
	HarvestAugmentPhysical  = HarvestType("AugmentPhysical")
	HarvestAugmentFire      = HarvestType("AugmentFire")
	HarvestAugmentAttack    = HarvestType("AugmentAttack")
	HarvestAugmentLife      = HarvestType("AugmentLife")
	HarvestAugmentCold      = HarvestType("AugmentCold")
	HarvestAugmentSpeed     = HarvestType("AugmentSpeed")
	HarvestAugmentDefence   = HarvestType("AugmentDefence")
	HarvestAugmentLightning = HarvestType("AugmentLightning")
	HarvestAugmentChaos     = HarvestType("AugmentChaos")
	HarvestAugmentCritical  = HarvestType("AugmentCritical")
	HarvestAugmentInfluence = HarvestType("AugmentInfluence")

	// Augment a * Magic or Rare item with a new * modifier with Lucky values

	HarvestAugmentCasterLucky    = HarvestType("AugmentCasterLucky")
	HarvestAugmentPhysicalLucky  = HarvestType("AugmentPhysicalLucky")
	HarvestAugmentFireLucky      = HarvestType("AugmentFireLucky")
	HarvestAugmentAttackLucky    = HarvestType("AugmentAttackLucky")
	HarvestAugmentLifeLucky      = HarvestType("AugmentLifeLucky")
	HarvestAugmentColdLucky      = HarvestType("AugmentColdLucky")
	HarvestAugmentSpeedLucky     = HarvestType("AugmentSpeedLucky")
	HarvestAugmentDefenceLucky   = HarvestType("AugmentDefenceLucky")
	HarvestAugmentLightningLucky = HarvestType("AugmentLightningLucky")
	HarvestAugmentChaosLucky     = HarvestType("AugmentChaosLucky")
	HarvestAugmentCriticalLucky  = HarvestType("AugmentCriticalLucky")
	HarvestAugmentInfluenceLucky = HarvestType("AugmentInfluenceLucky")

	// Remove a random Influence modifier from an item

	RemoveInfluence = HarvestType("Influence")

	// Reveals a random Currency crafting effect with improved outcome chances when Harvested

	HarvestReforgeSockets10 = HarvestType("ReforgeSockets10")
	HarvestReforgeLinks10   = HarvestType("ReforgeLinks10")
	HarvestReforgeColors10  = HarvestType("ReforgeColors10")
	HarvestQualityFlask     = HarvestType("QualityFlask")
	HarvestQualityGem       = HarvestType("QualityGem")
	HarvestUpgradeRare10    = HarvestType("UpgradeRare10")
	HarvestReforgeRare10    = HarvestType("ReforgeRare10")
	HarvestCorrupt10        = HarvestType("Corrupt10")

	// Reveals a random crafting effect that changes the element of Elemental Resistance modifiers when Harvested

	HarvestChangeResistColdToFire      = HarvestType("ChangeResistColdToFire")
	HarvestChangeResistColdToLightning = HarvestType("ChangeResistColdToLightning")
	HarvestChangeResistFireToCold      = HarvestType("ChangeResistFireToCold")
	HarvestChangeResistFireToLightning = HarvestType("ChangeResistFireToLightning")
	HarvestChangeResistLightningToCold = HarvestType("ChangeResistLightningToCold")
	HarvestChangeResistLightningToFire = HarvestType("ChangeResistLightningToFire")

	// Reveals a currency item exchange, trading three of a certain type of currency for other currency items when Harvested

	HarvestExchangeChaosToDivine          = HarvestType("ExchangeChaosToDivine")
	HarvestExchangeTransmutationToAlchemy = HarvestType("ExchangeTransmutationToAlchemy")
	HarvestExchangeBlessedToAlteration    = HarvestType("ExchangeBlessedToAlteration")
	HarvestExchangeAlchemyToChisels       = HarvestType("ExchangeAlchemyToChisels")
	HarvestExchangeChromaticToGCP         = HarvestType("ExchangeChromaticToGCP")
	HarvestExchangeJewellerToFusing       = HarvestType("ExchangeJewellerToFusing")
	HarvestExchangeAugmentToRegal         = HarvestType("ExchangeAugmentToRegal")
	HarvestExchangeWisdomToChance         = HarvestType("ExchangeWisdomToChance")
	HarvestExchangeSextantsToPrime        = HarvestType("ExchangeSextantsToPrime")
	HarvestExchangePrimeToAwakened        = HarvestType("ExchangePrimeToAwakened")
	HarvestExchangeScouringToAnnulment    = HarvestType("ExchangeScouringToAnnulment")
	HarvestExchangeAlterationToChaos      = HarvestType("ExchangeAlterationToChaos")
	HarvestExchangeVaalToRegret           = HarvestType("ExchangeVaalToRegret")
	HarvestExchangeChiselsToVaal          = HarvestType("ExchangeChiselsToVaal")

	// Reveals a random effect that exchanges Fossils, Essences, Delirium Orbs, Oils or Catalysts when Harvested

	HarvestChangeFossils   = HarvestType("ChangeFossils")
	HarvestChangeEssences  = HarvestType("ChangeEssences")
	HarvestChangeDelirium  = HarvestType("ChangeDelirium")
	HarvestChangeOils      = HarvestType("ChangeOils")
	HarvestChangeCatalysts = HarvestType("ChangeCatalysts")

	// Reveals a random crafting effect that reforges a Rare item a certain way when Harvested

	HarvestReforgeKeepPrefixes = HarvestType("ReforgeKeepPrefixes")
	HarvestReforgeKeepSuffixes = HarvestType("ReforgeKeepSuffixes")
	HarvestReforgeLessLikely   = HarvestType("ReforgeLessLikely")
	HarvestReforgeMoreLikely   = HarvestType("ReforgeMoreLikely")

	// Allows you to sacrifice a Weapon or Armour to create Jewellery or Jewels when Harvested

	HarvestSacrificeForBelt   = HarvestType("SacrificeForBelt")
	HarvestSacrificeForRing   = HarvestType("SacrificeForRing")
	HarvestSacrificeForAmulet = HarvestType("SacrificeForAmulet")
	HarvestSacrificeForJewel  = HarvestType("SacrificeForJewel")

	// Allows you to Sacrifice a Map to create items for the Atlas when Harvested

	HarvestSacrificeMapWhiteYellowFragments = HarvestType("SacrificeMapWhiteYellowFragments")
	HarvestSacrificeMapRedFragments         = HarvestType("SacrificeMapRedFragments")
	HarvestSacrificeMapCurrency             = HarvestType("SacrificeMapCurrency")
	HarvestSacrificeMapScarab               = HarvestType("SacrificeMapScarab")
	HarvestSacrificeMapElder                = HarvestType("SacrificeMapElder")
	HarvestSacrificeMapShaper               = HarvestType("SacrificeMapShaper")
	HarvestSacrificeMapSynthesis            = HarvestType("SacrificeMapSynthesis")
	HarvestSacrificeMapConqueror            = HarvestType("SacrificeMapConqueror")

	// Reveals a random Weapon Enchantment that replaces Quality's effect when Harvested

	HarvestEnchantWeaponCritical        = HarvestType("EnchantWeaponCritical")
	HarvestEnchantWeaponAccuracy        = HarvestType("EnchantWeaponAccuracy")
	HarvestEnchantWeaponAttackSpeed     = HarvestType("EnchantWeaponAttackSpeed")
	HarvestEnchantWeaponRange           = HarvestType("EnchantWeaponRange")
	HarvestEnchantWeaponElementalDamage = HarvestType("EnchantWeaponElementalDamage")
	HarvestEnchantWeaponAoE             = HarvestType("EnchantWeaponAoE")

	// Reveal a random crafting effect that locks a random modifier on an item when Harvested

	HarvestFractureRandom       = HarvestType("FractureRandom")
	HarvestFractureRandomSuffix = HarvestType("FractureRandomSuffix")
	HarvestFractureRandomPrefix = HarvestType("FractureRandomPrefix")

	// Reveals a random Socket number crafting effect when Harvested

	HarvestThreeSockets = HarvestType("ThreeSockets")
	HarvestFourSockets  = HarvestType("FourSockets")
	HarvestFiveSockets  = HarvestType("FiveSockets")
	HarvestSixSockets   = HarvestType("SixSockets")

	// Reveals a random Gem crafting effect when Harvested

	HarvestChangeGemToGem              = HarvestType("ChangeGemToGem")
	HarvestSacrificeCorruptedQuality20 = HarvestType("SacrificeCorruptedQuality20")
	HarvestSacrificeCorruptedQuality30 = HarvestType("SacrificeCorruptedQuality30")
	HarvestSacrificeCorruptedQuality40 = HarvestType("SacrificeCorruptedQuality40")
	HarvestSacrificeCorruptedQuality50 = HarvestType("SacrificeCorruptedQuality50")
	HarvestSacrificeCorruptedExp20     = HarvestType("SacrificeCorruptedExp20")
	HarvestSacrificeCorruptedExp30     = HarvestType("SacrificeCorruptedExp30")
	HarvestSacrificeCorruptedExp40     = HarvestType("SacrificeCorruptedExp40")
	HarvestSacrificeCorruptedExp50     = HarvestType("SacrificeCorruptedExp50")
	HarvestAttemptAwaken               = HarvestType("AttemptAwaken")

	// Allows you to change a Map into other Maps when Harvested

	HarvestChangeMapSameTier     = HarvestType("ChangeMapSameTier")
	HarvestSacrificeMapTierLower = HarvestType("SacrificeMapTierLower")
	HarvestSacrificeMapCorrupted = HarvestType("SacrificeMapCorrupted")

	// Allows you to add an Implicit modifier to certain Jewel types when Harvested

	HarvestSetImplicitCobaltCrimsonViridianPrismatic = HarvestType("SetImplicitCobaltCrimsonViridianPrismatic")
	HarvestSetImplicitAbyssTimeless                  = HarvestType("SetImplicitAbyssTimeless")
	HarvestSetImplicitCluster                        = HarvestType("SetImplicitCluster")

	// Reveals a random Flask Enchantment that depletes as it is used when Harvested

	HarvestEnchantFlaskDuration       = HarvestType("EnchantFlaskDuration")
	HarvestEnchantFlaskEffect         = HarvestType("EnchantFlaskEffect")
	HarvestEnchantFlaskMaxCharges     = HarvestType("EnchantFlaskMaxCharges")
	HarvestEnchantFlaskReducedCharges = HarvestType("EnchantFlaskReducedCharges")

	// Reveals a random Map Enchantment when Harvested

	HarvestEnchantMapRandomMod = HarvestType("EnchantMapRandomMod")
	HarvestEnchantMapSextant   = HarvestType("EnchantMapSextant")
	HarvestEnchantMapVaal      = HarvestType("EnchantMapVaal")
	HarvestEnchantMapSpirits   = HarvestType("EnchantMapSpirits")

	// Reveals a currency item exchange, trading ten of a certain type of currency for other currency items when Harvested

	HarvestExchange10ChaosToExalted       = HarvestType("Exchange10ChaosToExalted")
	HarvestExchange10TransmutationAlchemy = HarvestType("Exchange10TransmutationAlchemy")
	HarvestExchange10BlessedToAlteration  = HarvestType("Exchange10BlessedToAlteration")
	HarvestExchange10AlchemyToChisels     = HarvestType("Exchange10AlchemyToChisels")
	HarvestExchange10ChromaticToGCP       = HarvestType("Exchange10ChromaticToGCP")
	HarvestExchange10JewellerToFusing     = HarvestType("Exchange10JewellerToFusing")
	HarvestExchange10AugmentationToRegal  = HarvestType("Exchange10AugmentationToRegal")
	HarvestExchange10WisdomToChance       = HarvestType("Exchange10WisdomToChance")
	HarvestExchange10SextantsToPrime      = HarvestType("Exchange10SextantsToPrime")
	HarvestExchange10PrimeToAwakened      = HarvestType("Exchange10PrimeToAwakened")
	HarvestExchange10ScouringToAnnulment  = HarvestType("Exchange10ScouringToAnnulment")
	HarvestExchange10AlterationToChaos    = HarvestType("Exchange10AlterationToChaos")
	HarvestExchange10VaalToRegret         = HarvestType("Exchange10VaalToRegret")
	HarvestExchange10ChiselsToVaal        = HarvestType("Exchange10ChiselsToVaal")

	// Allows you to enhance specialised Currency a certain way when Harvested

	HarvestUpgradeEssenceTier      = HarvestType("UpgradeEssenceTier")
	HarvestUpgradeOilTier          = HarvestType("UpgradeOilTier")
	HarvestUpgradeEngineers        = HarvestType("UpgradeEngineers")
	HarvestExchangeResonatorFossil = HarvestType("ExchangeResonatorFossil")

	// Allows you to modify a Scarab a certain way when Harvested

	HarvestChangeScarabSameRarity = HarvestType("ChangeScarabSameRarity")
	HarvestUpgradeScarab          = HarvestType("UpgradeScarab")
	HarvestSplitScarab            = HarvestType("SplitScarab")

	// Allows you to upgrade an Offering to the Goddess when Harvested

	HarvestChangeOfferingToTribute    = HarvestType("ChangeOfferingToTribute")
	HarvestChangeOfferingToGift       = HarvestType("ChangeOfferingToGift")
	HarvestChangeOfferingToDedication = HarvestType("ChangeOfferingToDedication")

	// Allows you to give an item a Synthesised implicit modifier when Harvested

	HarvestSynthesiseItem = HarvestType("SynthesiseItem")

	// Reveals a random Socket link crafting effect when Harvested

	HarvestThreeLinks = HarvestType("ThreeLinks")
	HarvestFourLinks  = HarvestType("FourLinks")
	HarvestFiveLinks  = HarvestType("FiveLinks")
	HarvestSixLinks   = HarvestType("SixLinks")

	// Reveals a random Unique Item transformation effect when Harvested

	HarvestChangeUniqueToUnique             = HarvestType("ChangeUniqueToUnique")
	HarvestChangeUniqueToWeaponQuiver       = HarvestType("ChangeUniqueToWeaponQuiver")
	HarvestChangeUniqueArmour               = HarvestType("ChangeUniqueArmour")
	HarvestChangeUniqueJewellery            = HarvestType("ChangeUniqueJewellery")
	HarvestChangeUniqueWeaponToWeapon       = HarvestType("ChangeUniqueWeaponToWeapon")
	HarvestChangeUniqueArmourToArmour       = HarvestType("ChangeUniqueArmourToArmour")
	HarvestChangeUniqueJewelleryToJewellery = HarvestType("ChangeUniqueJewelleryToJewellery")
	HarvestChangeUniqueJewelToJewel         = HarvestType("ChangeUniqueJewelToJewel")

	// Allows you to sacrifice Divination Cards to create Divination Cards when Harvested

	HarvestChangeDivToDiv     = HarvestType("ChangeDivToDiv")
	HarvestSacrificeDivGamble = HarvestType("SacrificeDivGamble")
	HarvestSacrificeDivStack  = HarvestType("SacrificeDivStack")

	// Allows you to exchange certain Map Fragments for another of the same type when Harvested

	HarvestChangeFragmentSacrificeMortal = HarvestType("ChangeFragmentSacrificeMortal")
	HarvestChangeFragmentShaper          = HarvestType("ChangeFragmentShaper")
	HarvestChangeFragmentElder           = HarvestType("ChangeFragmentElder")
	HarvestChangeFragmentConqueror       = HarvestType("ChangeFragmentConqueror")

	// Reveals a random crafting effect that upgrades a Normal or Magic item's Rarity when Harvested

	HarvestUpgradeMagicTwoMod       = HarvestType("UpgradeMagicTwoMod")
	HarvestUpgradeMagicThreeMod     = HarvestType("UpgradeMagicThreeMod")
	HarvestUpgradeMagicFourMod      = HarvestType("UpgradeMagicFourMod")
	HarvestUpgradeMagicTwoModHigh   = HarvestType("UpgradeMagicTwoModHigh")
	HarvestUpgradeMagicThreeModHigh = HarvestType("UpgradeMagicThreeModHigh")
	HarvestUpgradeMagicFourModHigh  = HarvestType("UpgradeMagicFourModHigh")
	HarvestUpgradeNormalOneModHigh  = HarvestType("UpgradeNormalOneModHigh")
	HarvestUpgradeNormalTwoModHigh  = HarvestType("UpgradeNormalTwoModHigh")

	// Allows you to modify an item, resulting in Lucky modifier values when Harvested

	HarvestReforgeKeepPrefixesLucky        = HarvestType("ReforgeKeepPrefixesLucky")
	HarvestReforgeKeepSuffixesLucky        = HarvestType("ReforgeKeepSuffixesLucky")
	HarvestRerollPrefixSuffixImplicitLucky = HarvestType("RerollPrefixSuffixImplicitLucky")
	HarvestAugmentLucky                    = HarvestType("AugmentLucky")
	HarvestRerollPrefixLucky               = HarvestType("RerollPrefixLucky")
	HarvestRerollSuffixLucky               = HarvestType("RerollSuffixLucky")

	// Allows you to exchange Splinters, Breachstones or Emblems for others of the same type when Harvested

	HarvestChangeBreach   = HarvestType("ChangeBreach")
	HarvestChangeTimeless = HarvestType("ChangeTimeless")

	// Allows you to exchange a Unique, Bestiary or Harbinger item for a related item when Harvested

	HarvestChangeUniqueBestiary  = HarvestType("ChangeUniqueBestiary")
	HarvestChangeUniqueHarbinger = HarvestType("ChangeUniqueHarbinger")

	// Allows you to randomise the Atlas Influence types on an Influenced item when Harvested

	HarvestRandomiseInfluenceWeapon    = HarvestType("RandomiseInfluenceWeapon")
	HarvestRandomiseInfluenceArmour    = HarvestType("RandomiseInfluenceArmour")
	HarvestRandomiseInfluenceJewellery = HarvestType("RandomiseInfluenceJewellery")

	// Reveals a random Body Armour Enchantment that replaces Quality's effect when Harvested

	HarvestEnchantArmourLife            = HarvestType("EnchantArmourLife")
	HarvestEnchantArmourMana            = HarvestType("EnchantArmourMana")
	HarvestEnchantArmourStrength        = HarvestType("EnchantArmourStrength")
	HarvestEnchantArmourDexterity       = HarvestType("EnchantArmourDexterity")
	HarvestEnchantArmourIntelligence    = HarvestType("EnchantArmourIntelligence")
	HarvestEnchantArmourFireResist      = HarvestType("EnchantArmourFireResist")
	HarvestEnchantArmourColdResist      = HarvestType("EnchantArmourColdResist")
	HarvestEnchantArmourLightningResist = HarvestType("EnchantArmourLightningResist")

	// Randomise the numeric values of the random Influence modifiers on a Magic or Rare item

	HarvestRandomiseNumericInfluence = HarvestType("RandomiseNumericInfluence")
)

type HarvestCraft struct {
	Type    HarvestType `json:"type"`
	Message string      `json:"message"`
	Pricing string      `json:"pricing"`
}

type CraftWithText struct {
	Text  string
	Craft HarvestCraft
}

var crafts = make(map[string]HarvestCraft)
var reverseCrafts = make(map[HarvestType]CraftWithText)
var reversePricing = make(map[string]CraftWithText)

func InitCrafts() {
	if err := json.Unmarshal(CraftsJSON, &crafts); err != nil {
		panic(err)
	}

	for text, craft := range crafts {
		reverseCrafts[craft.Type] = CraftWithText{
			Text:  text,
			Craft: craft,
		}
	}

	for text, craft := range crafts {
		if craft.Pricing != "" {
			reversePricing[craft.Pricing] = CraftWithText{
				Text:  text,
				Craft: craft,
			}
		}
	}
}

func GetCraftByText(text string) HarvestCraft {
	return crafts[text]
}

func GetCraftByType(craftType HarvestType) CraftWithText {
	return reverseCrafts[craftType]
}

func GetCraftByPricing(pricing string) *CraftWithText {
	if craft, ok := reversePricing[pricing]; ok {
		return &craft
	}
	return nil
}

func FindCraft(text string) string {
	clean := strings.Replace(text, "\n", " ", -1)

	closestDistance := float64(0)
	closestString := ""

	for s := range crafts {
		dist := textdistance.JaroWinklerDistance(clean, s)
		if dist > closestDistance {
			closestString = s
			closestDistance = dist
		}
	}

	return closestString
}

func AllCrafts() map[string]HarvestCraft {
	return crafts
}

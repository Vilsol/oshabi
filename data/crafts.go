package data

import (
	"math"
	"strings"

	"github.com/vilsol/oshabi/config"
	"github.com/vilsol/oshabi/types"

	"github.com/masatana/go-textdistance"
)

//goland:noinspection GoUnusedConst
const (
	// Reveals a random Socket colour crafting effect when Harvested

	HarvestReforgeNonRedToRed        = types.HarvestType("ReforgeNonRedToRed")
	HarvestReforgeNonBlueToBlue      = types.HarvestType("ReforgeNonBlueToBlue")
	HarvestReforgeNonGreenToGreen    = types.HarvestType("ReforgeNonGreenToGreen")
	HarvestReforgeTwoRandomRedBlue   = types.HarvestType("ReforgeTwoRandomRedBlue")
	HarvestReforgeTwoRandomRedGreen  = types.HarvestType("ReforgeTwoRandomRedGreen")
	HarvestReforgeTwoRandomBlueGreen = types.HarvestType("ReforgeTwoRandomBlueGreen")
	HarvestReforgeThreeRandomRGB     = types.HarvestType("ReforgeThreeRandomRGB")
	HarvestReforgeWhite              = types.HarvestType("ReforgeWhite")

	// Reforge a Normal, Magic or Rare item as a Rare item with random modifiers, including a * modifier

	HarvestReforgeCaster    = types.HarvestType("ReforgeCaster")
	HarvestReforgePhysical  = types.HarvestType("ReforgePhysical")
	HarvestReforgeFire      = types.HarvestType("ReforgeFire")
	HarvestReforgeAttack    = types.HarvestType("ReforgeAttack")
	HarvestReforgeLife      = types.HarvestType("ReforgeLife")
	HarvestReforgeCold      = types.HarvestType("ReforgeCold")
	HarvestReforgeSpeed     = types.HarvestType("ReforgeSpeed")
	HarvestReforgeDefence   = types.HarvestType("ReforgeDefence")
	HarvestReforgeLightning = types.HarvestType("ReforgeLightning")
	HarvestReforgeChaos     = types.HarvestType("ReforgeChaos")
	HarvestReforgeCritical  = types.HarvestType("ReforgeCritical")
	HarvestReforgeInfluence = types.HarvestType("ReforgeInfluence")

	// Reforge a Normal, Magic or Rare item as a Rare item with random modifiers, including a * modifier. * modifiers are more common

	HarvestReforgeCasterMoreLikely    = types.HarvestType("ReforgeCasterMoreLikely")
	HarvestReforgePhysicalMoreLikely  = types.HarvestType("ReforgePhysicalMoreLikely")
	HarvestReforgeFireMoreLikely      = types.HarvestType("ReforgeFireMoreLikely")
	HarvestReforgeAttackMoreLikely    = types.HarvestType("ReforgeAttackMoreLikely")
	HarvestReforgeLifeMoreLikely      = types.HarvestType("ReforgeLifeMoreLikely")
	HarvestReforgeColdMoreLikely      = types.HarvestType("ReforgeColdMoreLikely")
	HarvestReforgeSpeedMoreLikely     = types.HarvestType("ReforgeSpeedMoreLikely")
	HarvestReforgeDefenceMoreLikely   = types.HarvestType("ReforgeDefenceMoreLikely")
	HarvestReforgeLightningMoreLikely = types.HarvestType("ReforgeLightningMoreLikely")
	HarvestReforgeChaosMoreLikely     = types.HarvestType("ReforgeChaosMoreLikely")
	HarvestReforgeCriticalMoreLikely  = types.HarvestType("ReforgeCriticalMoreLikely")
	HarvestReforgeInfluenceMoreLikely = types.HarvestType("ReforgeInfluenceMoreLikely")

	// Remove a random non-* modifier from a * item and add a new * modifier

	HarvestRemoveAddNonCaster    = types.HarvestType("RemoveAddNonCaster")
	HarvestRemoveAddNonPhysical  = types.HarvestType("RemoveAddNonPhysical")
	HarvestRemoveAddNonFire      = types.HarvestType("RemoveAddNonFire")
	HarvestRemoveAddNonAttack    = types.HarvestType("RemoveAddNonAttack")
	HarvestRemoveAddNonLife      = types.HarvestType("RemoveAddNonLife")
	HarvestRemoveAddNonCold      = types.HarvestType("RemoveAddNonCold")
	HarvestRemoveAddNonSpeed     = types.HarvestType("RemoveAddNonSpeed")
	HarvestRemoveAddNonDefence   = types.HarvestType("RemoveAddNonDefence")
	HarvestRemoveAddNonLightning = types.HarvestType("RemoveAddNonLightning")
	HarvestRemoveAddNonChaos     = types.HarvestType("RemoveAddNonChaos")
	HarvestRemoveAddNonCritical  = types.HarvestType("RemoveAddNonCritical")
	HarvestRemoveAddNonInfluence = types.HarvestType("RemoveAddNonInfluence")
	HarvestRemoveAddInfluence    = types.HarvestType("RemoveAddInfluence")

	// Augment a * Magic or Rare item with a new * modifier

	HarvestAugmentCaster    = types.HarvestType("AugmentCaster")
	HarvestAugmentPhysical  = types.HarvestType("AugmentPhysical")
	HarvestAugmentFire      = types.HarvestType("AugmentFire")
	HarvestAugmentAttack    = types.HarvestType("AugmentAttack")
	HarvestAugmentLife      = types.HarvestType("AugmentLife")
	HarvestAugmentCold      = types.HarvestType("AugmentCold")
	HarvestAugmentSpeed     = types.HarvestType("AugmentSpeed")
	HarvestAugmentDefence   = types.HarvestType("AugmentDefence")
	HarvestAugmentLightning = types.HarvestType("AugmentLightning")
	HarvestAugmentChaos     = types.HarvestType("AugmentChaos")
	HarvestAugmentCritical  = types.HarvestType("AugmentCritical")
	HarvestAugmentInfluence = types.HarvestType("AugmentInfluence")

	// Augment a * Magic or Rare item with a new * modifier with Lucky values

	HarvestAugmentCasterLucky    = types.HarvestType("AugmentCasterLucky")
	HarvestAugmentPhysicalLucky  = types.HarvestType("AugmentPhysicalLucky")
	HarvestAugmentFireLucky      = types.HarvestType("AugmentFireLucky")
	HarvestAugmentAttackLucky    = types.HarvestType("AugmentAttackLucky")
	HarvestAugmentLifeLucky      = types.HarvestType("AugmentLifeLucky")
	HarvestAugmentColdLucky      = types.HarvestType("AugmentColdLucky")
	HarvestAugmentSpeedLucky     = types.HarvestType("AugmentSpeedLucky")
	HarvestAugmentDefenceLucky   = types.HarvestType("AugmentDefenceLucky")
	HarvestAugmentLightningLucky = types.HarvestType("AugmentLightningLucky")
	HarvestAugmentChaosLucky     = types.HarvestType("AugmentChaosLucky")
	HarvestAugmentCriticalLucky  = types.HarvestType("AugmentCriticalLucky")
	HarvestAugmentInfluenceLucky = types.HarvestType("AugmentInfluenceLucky")

	// Remove a random Influence modifier from an item

	RemoveInfluence = types.HarvestType("Influence")

	// Reveals a random Currency crafting effect with improved outcome chances when Harvested

	HarvestReforgeSockets10 = types.HarvestType("ReforgeSockets10")
	HarvestReforgeLinks10   = types.HarvestType("ReforgeLinks10")
	HarvestReforgeColors10  = types.HarvestType("ReforgeColors10")
	HarvestQualityFlask     = types.HarvestType("QualityFlask")
	HarvestQualityGem       = types.HarvestType("QualityGem")
	HarvestUpgradeRare10    = types.HarvestType("UpgradeRare10")
	HarvestReforgeRare10    = types.HarvestType("ReforgeRare10")
	HarvestCorrupt10        = types.HarvestType("Corrupt10")

	// Reveals a random crafting effect that changes the element of Elemental Resistance modifiers when Harvested

	HarvestChangeResistColdToFire      = types.HarvestType("ChangeResistColdToFire")
	HarvestChangeResistColdToLightning = types.HarvestType("ChangeResistColdToLightning")
	HarvestChangeResistFireToCold      = types.HarvestType("ChangeResistFireToCold")
	HarvestChangeResistFireToLightning = types.HarvestType("ChangeResistFireToLightning")
	HarvestChangeResistLightningToCold = types.HarvestType("ChangeResistLightningToCold")
	HarvestChangeResistLightningToFire = types.HarvestType("ChangeResistLightningToFire")

	// Reveals a currency item exchange, trading three of a certain type of currency for other currency items when Harvested

	HarvestExchangeChaosToDivine          = types.HarvestType("ExchangeChaosToDivine")
	HarvestExchangeTransmutationToAlchemy = types.HarvestType("ExchangeTransmutationToAlchemy")
	HarvestExchangeBlessedToAlteration    = types.HarvestType("ExchangeBlessedToAlteration")
	HarvestExchangeAlchemyToChisels       = types.HarvestType("ExchangeAlchemyToChisels")
	HarvestExchangeChromaticToGCP         = types.HarvestType("ExchangeChromaticToGCP")
	HarvestExchangeJewellerToFusing       = types.HarvestType("ExchangeJewellerToFusing")
	HarvestExchangeAugmentToRegal         = types.HarvestType("ExchangeAugmentToRegal")
	HarvestExchangeWisdomToChance         = types.HarvestType("ExchangeWisdomToChance")
	HarvestExchangeSextantsToPrime        = types.HarvestType("ExchangeSextantsToPrime")
	HarvestExchangePrimeToAwakened        = types.HarvestType("ExchangePrimeToAwakened")
	HarvestExchangeScouringToAnnulment    = types.HarvestType("ExchangeScouringToAnnulment")
	HarvestExchangeAlterationToChaos      = types.HarvestType("ExchangeAlterationToChaos")
	HarvestExchangeVaalToRegret           = types.HarvestType("ExchangeVaalToRegret")
	HarvestExchangeChiselsToVaal          = types.HarvestType("ExchangeChiselsToVaal")

	// Reveals a random effect that exchanges Fossils, Essences, Delirium Orbs, Oils or Catalysts when Harvested

	HarvestChangeFossils   = types.HarvestType("ChangeFossils")
	HarvestChangeEssences  = types.HarvestType("ChangeEssences")
	HarvestChangeDelirium  = types.HarvestType("ChangeDelirium")
	HarvestChangeOils      = types.HarvestType("ChangeOils")
	HarvestChangeCatalysts = types.HarvestType("ChangeCatalysts")

	// Reveals a random crafting effect that reforges a Rare item a certain way when Harvested

	HarvestReforgeKeepPrefixes = types.HarvestType("ReforgeKeepPrefixes")
	HarvestReforgeKeepSuffixes = types.HarvestType("ReforgeKeepSuffixes")
	HarvestReforgeLessLikely   = types.HarvestType("ReforgeLessLikely")
	HarvestReforgeMoreLikely   = types.HarvestType("ReforgeMoreLikely")

	// Allows you to sacrifice a Weapon or Armour to create Jewellery or Jewels when Harvested

	HarvestSacrificeForBelt   = types.HarvestType("SacrificeForBelt")
	HarvestSacrificeForRing   = types.HarvestType("SacrificeForRing")
	HarvestSacrificeForAmulet = types.HarvestType("SacrificeForAmulet")
	HarvestSacrificeForJewel  = types.HarvestType("SacrificeForJewel")

	// Allows you to Sacrifice a Map to create items for the Atlas when Harvested

	HarvestSacrificeMapWhiteYellowFragments = types.HarvestType("SacrificeMapWhiteYellowFragments")
	HarvestSacrificeMapRedFragments         = types.HarvestType("SacrificeMapRedFragments")
	HarvestSacrificeMapCurrency             = types.HarvestType("SacrificeMapCurrency")
	HarvestSacrificeMapScarab               = types.HarvestType("SacrificeMapScarab")
	HarvestSacrificeMapElder                = types.HarvestType("SacrificeMapElder")
	HarvestSacrificeMapShaper               = types.HarvestType("SacrificeMapShaper")
	HarvestSacrificeMapSynthesis            = types.HarvestType("SacrificeMapSynthesis")
	HarvestSacrificeMapConqueror            = types.HarvestType("SacrificeMapConqueror")

	// Reveals a random Weapon Enchantment that replaces Quality's effect when Harvested

	HarvestEnchantWeaponCritical        = types.HarvestType("EnchantWeaponCritical")
	HarvestEnchantWeaponAccuracy        = types.HarvestType("EnchantWeaponAccuracy")
	HarvestEnchantWeaponAttackSpeed     = types.HarvestType("EnchantWeaponAttackSpeed")
	HarvestEnchantWeaponRange           = types.HarvestType("EnchantWeaponRange")
	HarvestEnchantWeaponElementalDamage = types.HarvestType("EnchantWeaponElementalDamage")
	HarvestEnchantWeaponAoE             = types.HarvestType("EnchantWeaponAoE")

	// Reveal a random crafting effect that locks a random modifier on an item when Harvested

	HarvestFractureRandom       = types.HarvestType("FractureRandom")
	HarvestFractureRandomSuffix = types.HarvestType("FractureRandomSuffix")
	HarvestFractureRandomPrefix = types.HarvestType("FractureRandomPrefix")

	// Reveals a random Socket number crafting effect when Harvested

	HarvestThreeSockets = types.HarvestType("ThreeSockets")
	HarvestFourSockets  = types.HarvestType("FourSockets")
	HarvestFiveSockets  = types.HarvestType("FiveSockets")
	HarvestSixSockets   = types.HarvestType("SixSockets")

	// Reveals a random Gem crafting effect when Harvested

	HarvestChangeGemToGem              = types.HarvestType("ChangeGemToGem")
	HarvestSacrificeCorruptedQuality20 = types.HarvestType("SacrificeCorruptedQuality20")
	HarvestSacrificeCorruptedQuality30 = types.HarvestType("SacrificeCorruptedQuality30")
	HarvestSacrificeCorruptedQuality40 = types.HarvestType("SacrificeCorruptedQuality40")
	HarvestSacrificeCorruptedQuality50 = types.HarvestType("SacrificeCorruptedQuality50")
	HarvestSacrificeCorruptedExp20     = types.HarvestType("SacrificeCorruptedExp20")
	HarvestSacrificeCorruptedExp30     = types.HarvestType("SacrificeCorruptedExp30")
	HarvestSacrificeCorruptedExp40     = types.HarvestType("SacrificeCorruptedExp40")
	HarvestSacrificeCorruptedExp50     = types.HarvestType("SacrificeCorruptedExp50")
	HarvestAttemptAwaken               = types.HarvestType("AttemptAwaken")

	// Allows you to change a Map into other Maps when Harvested

	HarvestChangeMapSameTier     = types.HarvestType("ChangeMapSameTier")
	HarvestSacrificeMapTierLower = types.HarvestType("SacrificeMapTierLower")
	HarvestSacrificeMapCorrupted = types.HarvestType("SacrificeMapCorrupted")

	// Allows you to add an Implicit modifier to certain Jewel types when Harvested

	HarvestSetImplicitCobaltCrimsonViridianPrismatic = types.HarvestType("SetImplicitCobaltCrimsonViridianPrismatic")
	HarvestSetImplicitAbyssTimeless                  = types.HarvestType("SetImplicitAbyssTimeless")
	HarvestSetImplicitCluster                        = types.HarvestType("SetImplicitCluster")

	// Reveals a random Flask Enchantment that depletes as it is used when Harvested

	HarvestEnchantFlaskDuration       = types.HarvestType("EnchantFlaskDuration")
	HarvestEnchantFlaskEffect         = types.HarvestType("EnchantFlaskEffect")
	HarvestEnchantFlaskMaxCharges     = types.HarvestType("EnchantFlaskMaxCharges")
	HarvestEnchantFlaskReducedCharges = types.HarvestType("EnchantFlaskReducedCharges")

	// Reveals a random Map Enchantment when Harvested

	HarvestEnchantMapRandomMod = types.HarvestType("EnchantMapRandomMod")
	HarvestEnchantMapSextant   = types.HarvestType("EnchantMapSextant")
	HarvestEnchantMapVaal      = types.HarvestType("EnchantMapVaal")
	HarvestEnchantMapSpirits   = types.HarvestType("EnchantMapSpirits")

	// Reveals a currency item exchange, trading ten of a certain type of currency for other currency items when Harvested

	HarvestExchange10ChaosToExalted       = types.HarvestType("Exchange10ChaosToExalted")
	HarvestExchange10TransmutationAlchemy = types.HarvestType("Exchange10TransmutationAlchemy")
	HarvestExchange10BlessedToAlteration  = types.HarvestType("Exchange10BlessedToAlteration")
	HarvestExchange10AlchemyToChisels     = types.HarvestType("Exchange10AlchemyToChisels")
	HarvestExchange10ChromaticToGCP       = types.HarvestType("Exchange10ChromaticToGCP")
	HarvestExchange10JewellerToFusing     = types.HarvestType("Exchange10JewellerToFusing")
	HarvestExchange10AugmentationToRegal  = types.HarvestType("Exchange10AugmentationToRegal")
	HarvestExchange10WisdomToChance       = types.HarvestType("Exchange10WisdomToChance")
	HarvestExchange10SextantsToPrime      = types.HarvestType("Exchange10SextantsToPrime")
	HarvestExchange10PrimeToAwakened      = types.HarvestType("Exchange10PrimeToAwakened")
	HarvestExchange10ScouringToAnnulment  = types.HarvestType("Exchange10ScouringToAnnulment")
	HarvestExchange10AlterationToChaos    = types.HarvestType("Exchange10AlterationToChaos")
	HarvestExchange10VaalToRegret         = types.HarvestType("Exchange10VaalToRegret")
	HarvestExchange10ChiselsToVaal        = types.HarvestType("Exchange10ChiselsToVaal")

	// Allows you to enhance specialised Currency a certain way when Harvested

	HarvestUpgradeEssenceTier      = types.HarvestType("UpgradeEssenceTier")
	HarvestUpgradeOilTier          = types.HarvestType("UpgradeOilTier")
	HarvestUpgradeEngineers        = types.HarvestType("UpgradeEngineers")
	HarvestExchangeResonatorFossil = types.HarvestType("ExchangeResonatorFossil")

	// Allows you to modify a Scarab a certain way when Harvested

	HarvestChangeScarabSameRarity = types.HarvestType("ChangeScarabSameRarity")
	HarvestUpgradeScarab          = types.HarvestType("UpgradeScarab")
	HarvestSplitScarab            = types.HarvestType("SplitScarab")

	// Allows you to upgrade an Offering to the Goddess when Harvested

	HarvestChangeOfferingToTribute    = types.HarvestType("ChangeOfferingToTribute")
	HarvestChangeOfferingToGift       = types.HarvestType("ChangeOfferingToGift")
	HarvestChangeOfferingToDedication = types.HarvestType("ChangeOfferingToDedication")

	// Allows you to give an item a Synthesised implicit modifier when Harvested

	HarvestSynthesiseItem = types.HarvestType("SynthesiseItem")

	// Reveals a random Socket link crafting effect when Harvested

	HarvestThreeLinks = types.HarvestType("ThreeLinks")
	HarvestFourLinks  = types.HarvestType("FourLinks")
	HarvestFiveLinks  = types.HarvestType("FiveLinks")
	HarvestSixLinks   = types.HarvestType("SixLinks")

	// Reveals a random Unique Item transformation effect when Harvested

	HarvestChangeUniqueToUnique             = types.HarvestType("ChangeUniqueToUnique")
	HarvestChangeUniqueToWeaponQuiver       = types.HarvestType("ChangeUniqueToWeaponQuiver")
	HarvestChangeUniqueArmour               = types.HarvestType("ChangeUniqueArmour")
	HarvestChangeUniqueJewellery            = types.HarvestType("ChangeUniqueJewellery")
	HarvestChangeUniqueWeaponToWeapon       = types.HarvestType("ChangeUniqueWeaponToWeapon")
	HarvestChangeUniqueArmourToArmour       = types.HarvestType("ChangeUniqueArmourToArmour")
	HarvestChangeUniqueJewelleryToJewellery = types.HarvestType("ChangeUniqueJewelleryToJewellery")
	HarvestChangeUniqueJewelToJewel         = types.HarvestType("ChangeUniqueJewelToJewel")

	// Allows you to sacrifice Divination Cards to create Divination Cards when Harvested

	HarvestChangeDivToDiv     = types.HarvestType("ChangeDivToDiv")
	HarvestSacrificeDivGamble = types.HarvestType("SacrificeDivGamble")
	HarvestSacrificeDivStack  = types.HarvestType("SacrificeDivStack")

	// Allows you to exchange certain Map Fragments for another of the same type when Harvested

	HarvestChangeFragmentSacrificeMortal = types.HarvestType("ChangeFragmentSacrificeMortal")
	HarvestChangeFragmentShaper          = types.HarvestType("ChangeFragmentShaper")
	HarvestChangeFragmentElder           = types.HarvestType("ChangeFragmentElder")
	HarvestChangeFragmentConqueror       = types.HarvestType("ChangeFragmentConqueror")

	// Reveals a random crafting effect that upgrades a Normal or Magic item's Rarity when Harvested

	HarvestUpgradeMagicTwoMod       = types.HarvestType("UpgradeMagicTwoMod")
	HarvestUpgradeMagicThreeMod     = types.HarvestType("UpgradeMagicThreeMod")
	HarvestUpgradeMagicFourMod      = types.HarvestType("UpgradeMagicFourMod")
	HarvestUpgradeMagicTwoModHigh   = types.HarvestType("UpgradeMagicTwoModHigh")
	HarvestUpgradeMagicThreeModHigh = types.HarvestType("UpgradeMagicThreeModHigh")
	HarvestUpgradeMagicFourModHigh  = types.HarvestType("UpgradeMagicFourModHigh")
	HarvestUpgradeNormalOneModHigh  = types.HarvestType("UpgradeNormalOneModHigh")
	HarvestUpgradeNormalTwoModHigh  = types.HarvestType("UpgradeNormalTwoModHigh")

	// Allows you to modify an item, resulting in Lucky modifier values when Harvested

	HarvestReforgeKeepPrefixesLucky        = types.HarvestType("ReforgeKeepPrefixesLucky")
	HarvestReforgeKeepSuffixesLucky        = types.HarvestType("ReforgeKeepSuffixesLucky")
	HarvestRerollPrefixSuffixImplicitLucky = types.HarvestType("RerollPrefixSuffixImplicitLucky")
	HarvestAugmentLucky                    = types.HarvestType("AugmentLucky")
	HarvestRerollPrefixLucky               = types.HarvestType("RerollPrefixLucky")
	HarvestRerollSuffixLucky               = types.HarvestType("RerollSuffixLucky")

	// Allows you to exchange Splinters, Breachstones or Emblems for others of the same type when Harvested

	HarvestChangeBreach   = types.HarvestType("ChangeBreach")
	HarvestChangeTimeless = types.HarvestType("ChangeTimeless")

	// Allows you to exchange a Unique, Bestiary or Harbinger item for a related item when Harvested

	HarvestChangeUniqueBestiary  = types.HarvestType("ChangeUniqueBestiary")
	HarvestChangeUniqueHarbinger = types.HarvestType("ChangeUniqueHarbinger")

	// Allows you to randomise the Atlas Influence types on an Influenced item when Harvested

	HarvestRandomiseInfluenceWeapon    = types.HarvestType("RandomiseInfluenceWeapon")
	HarvestRandomiseInfluenceArmour    = types.HarvestType("RandomiseInfluenceArmour")
	HarvestRandomiseInfluenceJewellery = types.HarvestType("RandomiseInfluenceJewellery")

	// Reveals a random Body Armour Enchantment that replaces Quality's effect when Harvested

	HarvestEnchantArmourLife            = types.HarvestType("EnchantArmourLife")
	HarvestEnchantArmourMana            = types.HarvestType("EnchantArmourMana")
	HarvestEnchantArmourStrength        = types.HarvestType("EnchantArmourStrength")
	HarvestEnchantArmourDexterity       = types.HarvestType("EnchantArmourDexterity")
	HarvestEnchantArmourIntelligence    = types.HarvestType("EnchantArmourIntelligence")
	HarvestEnchantArmourFireResist      = types.HarvestType("EnchantArmourFireResist")
	HarvestEnchantArmourColdResist      = types.HarvestType("EnchantArmourColdResist")
	HarvestEnchantArmourLightningResist = types.HarvestType("EnchantArmourLightningResist")

	// Randomise the numeric values of the random Influence modifiers on a Magic or Rare item

	HarvestRandomiseNumericInfluence = types.HarvestType("RandomiseNumericInfluence")
)

type HarvestCraft struct {
	Type         types.HarvestType          `json:"type"`
	Message      string                     `json:"message"`
	Pricing      string                     `json:"pricing"`
	Translations map[config.Language]string `json:"translations"`
	Short        map[config.Language]string `json:"short"`
}

var reversePricing = make(map[string]HarvestCraft)

func InitCrafts() {
	for _, craft := range crafts {
		if craft.Pricing != "" {
			reversePricing[craft.Pricing] = craft
		}
	}
}

func GetCraft(craftType types.HarvestType) HarvestCraft {
	return crafts[craftType]
}

func GetCraftByPricing(pricing string) *HarvestCraft {
	if craft, ok := reversePricing[pricing]; ok {
		return &craft
	}
	return nil
}

func FindCraft(text string) types.HarvestType {
	clean := strings.Replace(text, "\n", " ", -1)

	closestDistance := math.MaxInt
	closestString := HarvestReforgeNonRedToRed

	for craftType, craft := range crafts {
		dist := textdistance.DamerauLevenshteinDistance(clean, craft.Translations[config.Get().Language])
		if dist < closestDistance {
			closestString = craftType
			closestDistance = dist
		}
	}

	return closestString
}

func AllCrafts() map[types.HarvestType]HarvestCraft {
	return crafts
}

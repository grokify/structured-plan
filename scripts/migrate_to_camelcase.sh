#!/bin/bash
# migrate_to_camelcase.sh
# Migrates all snake_case JSON tags to camelCase in Go source files
#
# Usage: ./scripts/migrate_to_camelcase.sh
#
# This script should be run from the repository root directory.

set -e

echo "Migrating snake_case JSON tags to camelCase..."

# Function to convert snake_case to camelCase
snake_to_camel() {
    echo "$1" | sed -E 's/_([a-z])/\U\1/g'
}

# Find all Go files (excluding tests and vendor)
find . -name "*.go" -type f ! -name "*_test.go" ! -path "./vendor/*" | while read -r file; do
    # Apply all sed replacements
    sed -i '' \
        -e 's/json:"acceptance_criteria/json:"acceptanceCriteria/g' \
        -e 's/json:"access_control/json:"accessControl/g' \
        -e 's/json:"access_level/json:"accessLevel/g' \
        -e 's/json:"affected_segments/json:"affectedSegments/g' \
        -e 's/json:"affected_users/json:"affectedUsers/g' \
        -e 's/json:"age_range/json:"ageRange/g' \
        -e 's/json:"alert_threshold/json:"alertThreshold/g' \
        -e 's/json:"aligned_objectives/json:"alignedObjectives/g' \
        -e 's/json:"alternatives_considered/json:"alternativesConsidered/g' \
        -e 's/json:"api_specifications/json:"apiSpecifications/g' \
        -e 's/json:"appendix_refs/json:"appendixRefs/g' \
        -e 's/json:"approved_at/json:"approvedAt/g' \
        -e 's/json:"architecture_decisions/json:"architectureDecisions/g' \
        -e 's/json:"as_a/json:"asA/g' \
        -e 's/json:"at_rest/json:"atRest/g' \
        -e 's/json:"audit_logging/json:"auditLogging/g' \
        -e 's/json:"auth_method/json:"authMethod/g' \
        -e 's/json:"authentication_methods/json:"authenticationMethods/g' \
        -e 's/json:"authorization_model/json:"authorizationModel/g' \
        -e 's/json:"auto_scaling/json:"autoScaling/g' \
        -e 's/json:"backup_frequency/json:"backupFrequency/g' \
        -e 's/json:"base_url/json:"baseUrl/g' \
        -e 's/json:"below_threshold/json:"belowThreshold/g' \
        -e 's/json:"branch_strategy/json:"branchStrategy/g' \
        -e 's/json:"brand_guidelines/json:"brandGuidelines/g' \
        -e 's/json:"budget_authority/json:"budgetAuthority/g' \
        -e 's/json:"business_goals/json:"businessGoals/g' \
        -e 's/json:"buyer_personas/json:"buyerPersonas/g' \
        -e 's/json:"buying_criteria/json:"buyingCriteria/g' \
        -e 's/json:"buying_role/json:"buyingRole/g' \
        -e 's/json:"call_to_action/json:"callToAction/g' \
        -e 's/json:"category_scores/json:"categoryScores/g' \
        -e 's/json:"chosen_option_id/json:"chosenOptionId/g' \
        -e 's/json:"code_review/json:"codeReview/g' \
        -e 's/json:"coding_standards/json:"codingStandards/g' \
        -e 's/json:"collection_interval/json:"collectionInterval/g' \
        -e 's/json:"company_size/json:"companySize/g' \
        -e 's/json:"competitive_gaps/json:"competitiveGaps/g' \
        -e 's/json:"competitive_landscape/json:"competitiveLandscape/g' \
        -e 's/json:"compliance_controls/json:"complianceControls/g' \
        -e 's/json:"compliance_standards/json:"complianceStandards/g' \
        -e 's/json:"confidence_level/json:"confidenceLevel/g' \
        -e 's/json:"content_string/json:"contentString/g' \
        -e 's/json:"content_table/json:"contentTable/g' \
        -e 's/json:"correlation_id/json:"correlationId/g' \
        -e 's/json:"coverage_requirements/json:"coverageRequirements/g' \
        -e 's/json:"created_at/json:"createdAt/g' \
        -e 's/json:"current_alternatives/json:"currentAlternatives/g' \
        -e 's/json:"current_baseline/json:"currentBaseline/g' \
        -e 's/json:"current_state/json:"currentState/g' \
        -e 's/json:"current_value/json:"currentValue/g' \
        -e 's/json:"custom_metrics/json:"customMetrics/g' \
        -e 's/json:"custom_sections/json:"customSections/g' \
        -e 's/json:"customer_faqs/json:"customerFaqs/g' \
        -e 's/json:"customer_problem/json:"customerProblem/g' \
        -e 's/json:"customer_quote/json:"customerQuote/g' \
        -e 's/json:"data_classification/json:"dataClassification/g' \
        -e 's/json:"data_flows/json:"dataFlows/g' \
        -e 's/json:"data_format/json:"dataFormat/g' \
        -e 's/json:"data_model/json:"dataModel/g' \
        -e 's/json:"data_segregation/json:"dataSegregation/g' \
        -e 's/json:"data_stores/json:"dataStores/g' \
        -e 's/json:"data_type/json:"dataType/g' \
        -e 's/json:"ddos_protection/json:"ddosProtection/g' \
        -e 's/json:"decided_at/json:"decidedAt/g' \
        -e 's/json:"decided_by/json:"decidedBy/g' \
        -e 's/json:"design_principles/json:"designPrinciples/g' \
        -e 's/json:"design_system/json:"designSystem/g' \
        -e 's/json:"diagram_url/json:"diagramUrl/g' \
        -e 's/json:"disaster_recovery_plan/json:"disasterRecoveryPlan/g' \
        -e 's/json:"disaster_recovery/json:"disasterRecovery/g' \
        -e 's/json:"distribution_channels/json:"distributionChannels/g' \
        -e 's/json:"due_date/json:"dueDate/g' \
        -e 's/json:"encryption_at_rest/json:"encryptionAtRest/g' \
        -e 's/json:"encryption_in_transit/json:"encryptionInTransit/g' \
        -e 's/json:"encryption_model/json:"encryptionModel/g' \
        -e 's/json:"end_date/json:"endDate/g' \
        -e 's/json:"error_budget/json:"errorBudget/g' \
        -e 's/json:"escalation_policy/json:"escalationPolicy/g' \
        -e 's/json:"estimated_effort/json:"estimatedEffort/g' \
        -e 's/json:"executive_summary/json:"executiveSummary/g' \
        -e 's/json:"expected_outcomes/json:"expectedOutcomes/g' \
        -e 's/json:"export_format/json:"exportFormat/g' \
        -e 's/json:"failover_strategy/json:"failoverStrategy/g' \
        -e 's/json:"field_level/json:"fieldLevel/g' \
        -e 's/json:"future_phase/json:"futurePhase/g' \
        -e 's/json:"geographic_focus/json:"geographicFocus/g' \
        -e 's/json:"go_to_market/json:"goToMarket/g' \
        -e 's/json:"growth_rate/json:"growthRate/g' \
        -e 's/json:"high_availability/json:"highAvailability/g' \
        -e 's/json:"horizontal_scaling/json:"horizontalScaling/g' \
        -e 's/json:"how_it_works/json:"howItWorks/g' \
        -e 's/json:"i_want/json:"iWant/g' \
        -e 's/json:"image_url/json:"imageUrl/g' \
        -e 's/json:"in_scope/json:"inScope/g' \
        -e 's/json:"in_transit/json:"inTransit/g' \
        -e 's/json:"information_sources/json:"informationSources/g' \
        -e 's/json:"integration_points/json:"integrationPoints/g' \
        -e 's/json:"integration_tests/json:"integrationTests/g' \
        -e 's/json:"interaction_flows/json:"interactionFlows/g' \
        -e 's/json:"internal_faqs/json:"internalFaqs/g' \
        -e 's/json:"is_primary/json:"isPrimary/g' \
        -e 's/json:"isolation_model/json:"isolationModel/g' \
        -e 's/json:"issue_id/json:"issueId/g' \
        -e 's/json:"key_benefits/json:"keyBenefits/g' \
        -e 's/json:"key_decisions/json:"keyDecisions/g' \
        -e 's/json:"key_features/json:"keyFeatures/g' \
        -e 's/json:"key_findings/json:"keyFindings/g' \
        -e 's/json:"key_management/json:"keyManagement/g' \
        -e 's/json:"key_results/json:"keyResults/g' \
        -e 's/json:"key_threats/json:"keyThreats/g' \
        -e 's/json:"launch_strategy/json:"launchStrategy/g' \
        -e 's/json:"launch_timing/json:"launchTiming/g' \
        -e 's/json:"library_description/json:"libraryDescription/g' \
        -e 's/json:"library_ref/json:"libraryRef/g' \
        -e 's/json:"load_balancing/json:"loadBalancing/g' \
        -e 's/json:"log_levels/json:"logLevels/g' \
        -e 's/json:"made_by/json:"madeBy/g' \
        -e 's/json:"market_awareness/json:"marketAwareness/g' \
        -e 's/json:"market_opportunity/json:"marketOpportunity/g' \
        -e 's/json:"market_overview/json:"marketOverview/g' \
        -e 's/json:"market_position/json:"marketPosition/g' \
        -e 's/json:"market_requirements/json:"marketRequirements/g' \
        -e 's/json:"market_risks/json:"marketRisks/g' \
        -e 's/json:"market_share/json:"marketShare/g' \
        -e 's/json:"market_stage/json:"marketStage/g' \
        -e 's/json:"marketing_strategy/json:"marketingStrategy/g' \
        -e 's/json:"max_points/json:"maxPoints/g' \
        -e 's/json:"max_score/json:"maxScore/g' \
        -e 's/json:"measurement_method/json:"measurementMethod/g' \
        -e 's/json:"message_queues/json:"messageQueues/g' \
        -e 's/json:"metrics_quality/json:"metricsQuality/g' \
        -e 's/json:"multi_tenancy/json:"multiTenancy/g' \
        -e 's/json:"network_isolation/json:"networkIsolation/g' \
        -e 's/json:"network_policy/json:"networkPolicy/g' \
        -e 's/json:"network_security/json:"networkSecurity/g' \
        -e 's/json:"noisy_neighbor_protection/json:"noisyNeighborProtection/g' \
        -e 's/json:"non_functional/json:"nonFunctional/g' \
        -e 's/json:"non_goals/json:"nonGoals/g' \
        -e 's/json:"okr_ref/json:"okrRef/g' \
        -e 's/json:"on_call_integration/json:"onCallIntegration/g' \
        -e 's/json:"open_items/json:"openItems/g' \
        -e 's/json:"open_questions/json:"openQuestions/g' \
        -e 's/json:"optional_complete/json:"optionalComplete/g' \
        -e 's/json:"optional_total/json:"optionalTotal/g' \
        -e 's/json:"out_of_scope/json:"outOfScope/g' \
        -e 's/json:"overall_decision/json:"overallDecision/g' \
        -e 's/json:"overall_score/json:"overallScore/g' \
        -e 's/json:"pain_points/json:"painPoints/g' \
        -e 's/json:"partner_strategy/json:"partnerStrategy/g' \
        -e 's/json:"penetration_testing/json:"penetrationTesting/g' \
        -e 's/json:"performance_tests/json:"performanceTests/g' \
        -e 's/json:"persona_id/json:"personaId/g' \
        -e 's/json:"phase_id/json:"phaseId/g' \
        -e 's/json:"prd_id/json:"prdId/g' \
        -e 's/json:"preferred_channels/json:"preferredChannels/g' \
        -e 's/json:"press_release/json:"pressRelease/g' \
        -e 's/json:"pricing_strategy/json:"pricingStrategy/g' \
        -e 's/json:"primary_metric/json:"primaryMetric/g' \
        -e 's/json:"primary_segments/json:"primarySegments/g' \
        -e 's/json:"problem_definition/json:"problemDefinition/g' \
        -e 's/json:"problem_solved/json:"problemSolved/g' \
        -e 's/json:"problem_statement/json:"problemStatement/g' \
        -e 's/json:"problem_summary/json:"problemSummary/g' \
        -e 's/json:"problems_addressed/json:"problemsAddressed/g' \
        -e 's/json:"product_goals/json:"productGoals/g' \
        -e 's/json:"proof_points/json:"proofPoints/g' \
        -e 's/json:"propagation_format/json:"propagationFormat/g' \
        -e 's/json:"proposed_offering/json:"proposedOffering/g' \
        -e 's/json:"proposed_solution/json:"proposedSolution/g' \
        -e 's/json:"quality_scores/json:"qualityScores/g' \
        -e 's/json:"rate_limit/json:"rateLimit/g' \
        -e 's/json:"recommendation_rationale/json:"recommendationRationale/g' \
        -e 's/json:"recommendation_summary/json:"recommendationSummary/g' \
        -e 's/json:"recommended_owner/json:"recommendedOwner/g' \
        -e 's/json:"referenced_by/json:"referencedBy/g' \
        -e 's/json:"related_documents/json:"relatedDocuments/g' \
        -e 's/json:"related_ids/json:"relatedIds/g' \
        -e 's/json:"required_actions/json:"requiredActions/g' \
        -e 's/json:"required_complete/json:"requiredComplete/g' \
        -e 's/json:"required_total/json:"requiredTotal/g' \
        -e 's/json:"requirements_quality/json:"requirementsQuality/g' \
        -e 's/json:"retention_period/json:"retentionPeriod/g' \
        -e 's/json:"review_board_summary/json:"reviewBoardSummary/g' \
        -e 's/json:"revision_history/json:"revisionHistory/g' \
        -e 's/json:"revision_triggers/json:"revisionTriggers/g' \
        -e 's/json:"risk_management/json:"riskManagement/g' \
        -e 's/json:"root_causes/json:"rootCauses/g' \
        -e 's/json:"sales_strategy/json:"salesStrategy/g' \
        -e 's/json:"sample_size/json:"sampleSize/g' \
        -e 's/json:"sampling_rate/json:"samplingRate/g' \
        -e 's/json:"scalability_design/json:"scalabilityDesign/g' \
        -e 's/json:"schema_version/json:"schemaVersion/g' \
        -e 's/json:"scope_discipline/json:"scopeDiscipline/g' \
        -e 's/json:"secondary_metrics/json:"secondaryMetrics/g' \
        -e 's/json:"secondary_problems/json:"secondaryProblems/g' \
        -e 's/json:"secondary_segments/json:"secondarySegments/g' \
        -e 's/json:"security_audit_frequency/json:"securityAuditFrequency/g' \
        -e 's/json:"security_controls/json:"securityControls/g' \
        -e 's/json:"security_design/json:"securityDesign/g' \
        -e 's/json:"security_model/json:"securityModel/g' \
        -e 's/json:"security_tests/json:"securityTests/g' \
        -e 's/json:"selected_solution_id/json:"selectedSolutionId/g' \
        -e 's/json:"semantic_versioning/json:"semanticVersioning/g' \
        -e 's/json:"sensitive_data_handling/json:"sensitiveDataHandling/g' \
        -e 's/json:"session_management/json:"sessionManagement/g' \
        -e 's/json:"slo_target/json:"sloTarget/g' \
        -e 's/json:"so_that/json:"soThat/g' \
        -e 's/json:"solution_fit/json:"solutionFit/g' \
        -e 's/json:"solution_options/json:"solutionOptions/g' \
        -e 's/json:"solution_rationale/json:"solutionRationale/g' \
        -e 's/json:"spec_url/json:"specUrl/g' \
        -e 's/json:"start_date/json:"startDate/g' \
        -e 's/json:"story_points/json:"storyPoints/g' \
        -e 's/json:"success_criteria/json:"successCriteria/g' \
        -e 's/json:"success_metrics/json:"successMetrics/g' \
        -e 's/json:"system_diagram/json:"systemDiagram/g' \
        -e 's/json:"target_audience/json:"targetAudience/g' \
        -e 's/json:"target_buyer/json:"targetBuyer/g' \
        -e 's/json:"target_date/json:"targetDate/g' \
        -e 's/json:"target_market/json:"targetMarket/g' \
        -e 's/json:"target_state/json:"targetState/g' \
        -e 's/json:"target_uptime/json:"targetUptime/g' \
        -e 's/json:"target_value/json:"targetValue/g' \
        -e 's/json:"team_needs/json:"teamNeeds/g' \
        -e 's/json:"technical_approach/json:"technicalApproach/g' \
        -e 's/json:"technical_architecture/json:"technicalArchitecture/g' \
        -e 's/json:"technical_faqs/json:"technicalFaqs/g' \
        -e 's/json:"technical_feasibility/json:"technicalFeasibility/g' \
        -e 's/json:"technical_proficiency/json:"technicalProficiency/g' \
        -e 's/json:"technology_stack/json:"technologyStack/g' \
        -e 's/json:"test_environments/json:"testEnvironments/g' \
        -e 's/json:"testing_approach/json:"testingApproach/g' \
        -e 's/json:"threat_actors/json:"threatActors/g' \
        -e 's/json:"threat_level/json:"threatLevel/g' \
        -e 's/json:"threat_model/json:"threatModel/g' \
        -e 's/json:"top_risks/json:"topRisks/g' \
        -e 's/json:"trust_boundaries/json:"trustBoundaries/g' \
        -e 's/json:"unit_tests/json:"unitTests/g' \
        -e 's/json:"updated_at/json:"updatedAt/g' \
        -e 's/json:"used_in_prds/json:"usedInPrds/g' \
        -e 's/json:"user_impact/json:"userImpact/g' \
        -e 's/json:"user_stories/json:"userStories/g' \
        -e 's/json:"user_story_ids/json:"userStoryIds/g' \
        -e 's/json:"user_understanding/json:"userUnderstanding/g' \
        -e 's/json:"ux_coverage/json:"uxCoverage/g' \
        -e 's/json:"ux_requirements/json:"uxRequirements/g' \
        -e 's/json:"value_proposition/json:"valueProposition/g' \
        -e 's/json:"vertical_scaling/json:"verticalScaling/g' \
        -e 's/json:"vulnerability_scanning/json:"vulnerabilityScanning/g' \
        -e 's/json:"weighted_score/json:"weightedScore/g' \
        -e 's/json:"why_not_chosen/json:"whyNotChosen/g' \
        "$file"
done

echo "Migration complete!"
echo ""
echo "Next steps:"
echo "1. Run 'go test ./...' to verify tests pass"
echo "2. Run 'golangci-lint run' to check for issues"
echo "3. Update example JSON files in examples/"
echo "4. Regenerate JSON schemas"
